import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
import os
from matplotlib.colors import LinearSegmentedColormap
from matplotlib.patches import FancyBboxPatch, FancyArrowPatch

os.makedirs("out/haco", exist_ok=True)

df = pd.read_csv("haco.csv")

min_scores = {
    "tsp/st70.tsp": 675.000000,
    "tsp/ulysses22.tsp": 6981.326677,
    "tsp/swiss42.tsp": 1273,
    "tsp/eil51.tsp": 426.000000,
    "tsp/eil76.tsp": 538.000000,
    "tsp/berlin52.tsp": 7542,
}

def normalize_score(row):
    return row["score"] / min_scores[row["file"]]

def relative_error(row):
    return (row["score"] - min_scores[row["file"]]) / min_scores[row["file"]]

df["score_normalized"] = df.apply(normalize_score, axis=1)
df["relative_error"] = df.apply(relative_error, axis=1)

def parse_name(name):
    parts = name.split('_')
    if len(parts) == 3:
        mutation, selection, crossover = parts
        
        mutation_ru = {"uniform": "Равномерная", "gauss": "Гауссовская"}.get(mutation, mutation)
        selection_ru = {"best": "Элитарная", "tournament": "Турнирная", "roulette": "Рулеточная"}.get(selection, selection)
        crossover_ru = {"ariphmetic": "Арифметический", "blx": "BLX", "sbx": "SBX"}.get(crossover, crossover)
        
        return mutation_ru, selection_ru, crossover_ru, f"{mutation_ru} | {selection_ru} | {crossover_ru}"
    return name, "", "", name

df[['mutation', 'selection', 'crossover', 'name_ru']] = df['name'].apply(lambda x: pd.Series(parse_name(x)))

agg = df.groupby(["mutation", "selection", "crossover", "name_ru", "file"]).agg({
    "score_normalized": ["mean", "std"],
    "relative_error": ["mean", "std"],
    "score": ["mean", "std"]
}).reset_index()

agg.columns = ["mutation", "selection", "crossover", "name_ru", "file",
               "score_norm_mean", "score_norm_std",
               "rel_error_mean", "rel_error_std",
               "score_mean", "score_std"]

agg['rank'] = agg.groupby('file')['score_norm_mean'].rank(method='dense', ascending=True)
agg['file_clean'] = agg['file'].str.replace('tsp/', '', regex=False)

algorithm_stats = agg.groupby(['mutation', 'selection', 'crossover', 'name_ru']).agg({
    'rank': 'mean',
    'rel_error_mean': 'mean',
    'rel_error_std': 'mean'
}).reset_index()

algorithm_stats.columns = ['mutation', 'selection', 'crossover', 'name_ru', 'rank_mean', 'rel_error_mean', 'rel_error_std']
algorithm_stats = algorithm_stats.sort_values('rank_mean')
algorithm_stats.insert(0, '№', range(1, len(algorithm_stats) + 1))
algorithm_stats = algorithm_stats.rename(columns={
    'mutation': 'Мутация', 'selection': 'Селекция', 'crossover': 'Кроссинговер',
    'name_ru': 'Алгоритм', 'rank_mean': 'Средний ранг',
    'rel_error_mean': 'Ср. относительная ошибка', 'rel_error_std': 'СКО относительной ошибки'
})

print("\n=== Ранжирование алгоритмов ===")
print(algorithm_stats[['№', 'Мутация', 'Селекция', 'Кроссинговер', 'Средний ранг', 
                       'Ср. относительная ошибка', 'СКО относительной ошибки']].to_string(index=False, float_format='%.4f'))

latex_table = """\\begin{longtable}{|p{0.05\\textwidth}|p{0.15\\textwidth}|p{0.18\\textwidth}|p{0.15\\textwidth}|p{0.12\\textwidth}|p{0.06\\textwidth}|p{0.06\\textwidth}|}
\\caption{Ранжирование алгоритмов по среднему рангу}\\label{tbl:algorithm_ranking} \\\\
\\hline
\\textbf{№} & \\textbf{Селекция} & \\textbf{Кроссинговер} & \\textbf{Мутация} & \\textbf{Средний ранг} & \\textbf{$\\mu_{err}$} & \\textbf{$\\sigma_{err}$} \\\\
\\hline
\\endfirsthead
\\caption*{Продолжение таблицы~\\ref{tbl:algorithm_ranking}} \\\\
\\hline
\\textbf{№} & \\textbf{Селекция} & \\textbf{Кроссинговер} & \\textbf{Мутация} & \\textbf{Средний ранг} & \\textbf{$\\mu_{err}$} & \\textbf{$\\sigma_{err}$} \\\\
\\hline
\\endhead
\\hline
\\endfoot
\\endlastfoot
"""

for _, row in algorithm_stats.iterrows():
    latex_table += f"{int(row['№'])} & {row['Селекция']} & {row['Кроссинговер']} & {row['Мутация']} & {row['Средний ранг']:.1f} & {row['Ср. относительная ошибка']:.3f} & {row['СКО относительной ошибки']:.3f} \\\\\n\\hline\n"

latex_table += "\\end{longtable}\n"

with open("out/rank.tex", "w", encoding="utf-8") as f:
    f.write(latex_table)

print("\nLaTeX таблица сохранена в 'out/rank.tex'")

file_stats = agg.groupby(['file_clean', 'mutation', 'selection', 'crossover', 'name_ru']).agg({
    'score_mean': 'mean', 'score_std': 'mean', 'score_norm_mean': 'mean'
}).reset_index()

file_stats['rse'] = (file_stats['score_std'] / file_stats['score_mean']) * 100
file_stats = file_stats.sort_values(['file_clean', 'score_mean'])

variations_latex = ""
for file_name in file_stats['file_clean'].unique():
    file_data = file_stats[file_stats['file_clean'] == file_name].reset_index(drop=True)
    file_data.insert(0, '№', range(1, len(file_data) + 1))
    best_score = min_scores["tsp/" + file_name]
    
    variations_latex += f"\n% === Таблица для файла {file_name} ===\n"
    variations_latex += f"\\begin{{longtable}}{{|p{{0.05\\textwidth}}|p{{0.15\\textwidth}}|p{{0.18\\textwidth}}|p{{0.15\\textwidth}}|p{{0.12\\textwidth}}|p{{0.09\\textwidth}}|p{{0.1\\textwidth}}|}}\n"
    variations_latex += f"\\caption{{Результаты для задачи {file_name}. Оптимальное решение {best_score:.2f}}}\\label{{tbl:{file_name}}} \\\\\n\\hline\n"
    variations_latex += "\\textbf{№} & \\textbf{Селекция} & \\textbf{Кроссинговер} & \\textbf{Мутация} & \\textbf{Среднее} & \\textbf{СКО} & \\textbf{RSE (\\%)} \\\\\n\\hline\n"
    variations_latex += "\\endfirsthead\n"
    variations_latex += f"\\caption*{{Продолжение таблицы~\\ref{{tbl:{file_name}}}}} \\\\\n\\hline\n"
    variations_latex += "\\textbf{№} & \\textbf{Селекция} & \\textbf{Кроссинговер} & \\textbf{Мутация} & \\textbf{Среднее} & \\textbf{СКО} & \\textbf{RSE (\\%)} \\\\\n\\hline\n"
    variations_latex += "\\endhead\n\\hline\n\\endfoot\n\\endlastfoot\n"
    
    for _, row in file_data.iterrows():
        is_best = (row['score_mean'] == best_score)
        score_mean_str = f"\\textbf{{{row['score_mean']:.2f}}}" if is_best else f"{row['score_mean']:.2f}"
        variations_latex += f"{int(row['№'])} & {row['selection']} & {row['crossover']} & {row['mutation']} & {score_mean_str} & {row['score_std']:.2f} & {row['rse']:.2f} \\\\\n\\hline\n"
    
    variations_latex += "\\end{longtable}\n\n"

with open("out/variations.tex", "w", encoding="utf-8") as f:
    f.write(variations_latex)

print("LaTeX таблицы для каждого файла сохранены в 'out/variations.tex'")

# Нормальные размеры шрифтов для всех графиков
plt.rcParams.update({
    'font.size': 114,           # базовый размер шрифта
    'axes.labelsize': 116,      # размер подписей осей
    'axes.titlesize': 118,      # размер заголовка
    'xtick.labelsize': 112,     # размер подписей меток X
    'ytick.labelsize': 112,     # размер подписей меток Y
    'legend.fontsize': 112,     # размер шрифта легенды
    'figure.titlesize': 120     # размер заголовка фигуры
})

# График рангов с выделением лучшей конфигурации (золотой столбец + красный шрифт)
fig_rank, ax_rank = plt.subplots(figsize=(16, 10))

# Находим лучшую конфигурацию (с минимальным средним рангом)
best_idx = algorithm_stats['Средний ранг'].idxmin()
best_config = algorithm_stats.iloc[best_idx]
best_rank_value = best_config['Средний ранг']

# Создаем цвета: для лучшей - золотой, для остальных - из палитры
colors_list = []
for i in range(len(algorithm_stats)):
    if i == best_idx:
        colors_list.append('gold')
    else:
        colors_list.append(plt.cm.Set3(i / len(algorithm_stats)))

bars = ax_rank.bar(range(len(algorithm_stats)), algorithm_stats['Средний ранг'], 
                   color=colors_list, edgecolor='black', linewidth=1.5, zorder=2)

# Добавляем значения над столбцами
for i, (bar, val) in enumerate(zip(bars, algorithm_stats['Средний ранг'])):
    if i == best_idx:
        # Для лучшей конфигурации - красный жирный шрифт
        color = 'red'
        fontweight = 'bold'
        fontsize = 14
    else:
        color = 'black'
        fontweight = 'normal'
        fontsize = 13
    ax_rank.text(bar.get_x() + bar.get_width()/2, bar.get_height() + 0.1, 
                 f'{val:.2f}', ha='center', va='bottom', fontsize=fontsize, 
                 fontweight=fontweight, color=color)

# Добавляем стрелку к лучшему столбцу
best_bar = bars[best_idx]
arrow_y = best_bar.get_height() + 15
arrow = FancyArrowPatch((best_bar.get_x() + best_bar.get_width()/2, arrow_y),
                        (best_bar.get_x() + best_bar.get_width()/2, best_bar.get_height() + 0.2),
                        arrowstyle='->', mutation_scale=30, linewidth=2.5,
                        color='red', zorder=3)
ax_rank.add_patch(arrow)

# Добавляем текст со стрелкой
ax_rank.annotate('Лучшая конфигурация',
                xy=(best_bar.get_x() + best_bar.get_width()/2, best_bar.get_height() + 15.6),
                xytext=(best_bar.get_x() + best_bar.get_width()/2 + 1.5, best_bar.get_height() + 15.6),
                # arrowprops=dict(arrowstyle='->', color='red', lw=2, 
                #                connectionstyle='arc3,rad=0.3'),
                fontsize=12, fontweight='bold', color='red',
                ha='center', va='bottom',
                bbox=dict(boxstyle='round,pad=0.3', facecolor='lightyellow', 
                         edgecolor='red', alpha=0.9))

# Добавляем рамку вокруг лучшего столбца
rect = FancyBboxPatch((best_bar.get_x() - 0.03, best_bar.get_y() - 0.03),
                       best_bar.get_width() + 0.06, best_bar.get_height() + 0.06,
                       boxstyle='round,pad=0.02', linewidth=2.5, 
                       edgecolor='red', facecolor='none', zorder=4)
ax_rank.add_patch(rect)

# Делаем подпись оси X для лучшей конфигурации красной
xticklabels = [f'{label.get_text()}' for label in ax_rank.get_xticklabels()]
ax_rank.set_xlabel('Конфигурация', fontsize=16, fontweight='bold')
ax_rank.set_ylabel('Средний ранг', fontsize=16, fontweight='bold')
ax_rank.set_title('Средний ранг конфигураций', 
                  fontsize=18, fontweight='bold', pad=20)
ax_rank.set_xticks(range(len(algorithm_stats)))
ax_rank.set_xticklabels(algorithm_stats['Алгоритм'], rotation=45, ha='right', fontsize=11)

# Делаем красной метку лучшей конфигурации на оси X
labels = ax_rank.get_xticklabels()
labels[best_idx].set_color('red')
labels[best_idx].set_fontweight('bold')
labels[best_idx].set_fontsize(12)

ax_rank.grid(True, axis='y', linestyle='--', alpha=0.6, linewidth=0.8, zorder=1)
ax_rank.set_axisbelow(True)
ax_rank.tick_params(axis='both', which='major', labelsize=12)

# Устанавливаем отступ сверху для размещения аннотации
ax_rank.set_ylim(0, ax_rank.get_ylim()[1] * 1.15)

plt.tight_layout()
plt.savefig('out/haco/rank.png', dpi=200, bbox_inches='tight', facecolor='white')
plt.show()

# Тепловая карта с выделением лучшей конфигурации (красный шрифт)
rank_pivot = agg.pivot_table(index='name_ru', columns='file_clean', values='rank', aggfunc='first')
order = algorithm_stats.sort_values('Средний ранг')['Алгоритм'].tolist()
rank_pivot = rank_pivot.reindex(order)

fig_heatmap, ax_heatmap = plt.subplots(figsize=(14, 12))

# Создаем маску для выделения строки лучшей конфигурации
best_algorithm_name = best_config['Алгоритм']
best_row_idx = rank_pivot.index.get_loc(best_algorithm_name)

# Создаем пользовательскую цветовую карту
cmap = LinearSegmentedColormap.from_list('RdYlGn', ['green', 'yellow', 'red'], N=100)

im = ax_heatmap.imshow(rank_pivot.values, cmap=cmap, aspect='auto', 
                        interpolation='nearest', vmin=1, vmax=rank_pivot.values.max())

# Добавляем значения в ячейки
for i in range(rank_pivot.shape[0]):
    for j in range(rank_pivot.shape[1]):
        value = rank_pivot.values[i, j]
        if not pd.isna(value):
            norm_val = (value - 1) / (rank_pivot.values.max() - 1) if rank_pivot.values.max() > 1 else 0.5
            text_color = 'black'
            # Для лучшей конфигурации - красный жирный шрифт независимо от фона
            if i == best_row_idx:
                fontweight = 'bold'
                fontsize = 13
            else:
                fontweight = 'normal'
                fontsize = 12
            ax_heatmap.text(j, i, f'{value:.0f}', ha='center', va='center', 
                          fontsize=fontsize, fontweight=fontweight, color=text_color)

# Выделяем строку лучшей конфигурации красной рамкой
for j in range(rank_pivot.shape[1]):
    rect = plt.Rectangle((j - 0.5, best_row_idx - 0.5), 1, 1, 
                         fill=False, edgecolor='red', linewidth=3, 
                         linestyle='-', zorder=5)
    ax_heatmap.add_patch(rect)

# Добавляем звездочку и подпись к лучшей конфигурации
ax_heatmap.text(rank_pivot.shape[1] - 0.5, best_row_idx, ' ★', 
                fontsize=20, color='red', fontweight='bold', 
                ha='left', va='center', zorder=6)

# Делаем название лучшей конфигурации на оси Y красным
y_labels = [label.get_text() for label in ax_heatmap.get_yticklabels()]
ax_heatmap.set_yticks(range(len(rank_pivot.index)))
ax_heatmap.set_yticklabels(rank_pivot.index, fontsize=11)
yticklabels = ax_heatmap.get_yticklabels()
yticklabels[best_row_idx].set_color('red')
yticklabels[best_row_idx].set_fontweight('bold')
yticklabels[best_row_idx].set_fontsize(13)

ax_heatmap.set_xticks(range(len(rank_pivot.columns)))
ax_heatmap.set_xticklabels(rank_pivot.columns, rotation=45, ha='right', fontsize=13)
ax_heatmap.set_xlabel('Задача TSP', fontsize=16, fontweight='bold')
ax_heatmap.set_ylabel('Конфигурация', fontsize=16, fontweight='bold')
ax_heatmap.set_title('Тепловая карта рангов алгоритмов по задачам', 
                     fontsize=16, fontweight='bold', pad=20)

cbar = plt.colorbar(im, ax=ax_heatmap)
cbar.set_label('Ранг', fontsize=14, fontweight='bold')
cbar.ax.tick_params(labelsize=11)

plt.tight_layout()
plt.savefig('out/haco/heatmap.png', dpi=200, bbox_inches='tight', facecolor='white')
plt.show()

print(f"\nЛучшая конфигурация: {best_config['Алгоритм']} со средним рангом {best_config['Средний ранг']:.2f}")
print("\nГрафики сохранены:")
print("  - out/haco/rank.png")
print("  - out/haco/heatmap.png")
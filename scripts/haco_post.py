import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
import os
from matplotlib.colors import LinearSegmentedColormap

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

plt.rcParams.update({'font.size': 12, 'axes.labelsize': 14, 'axes.titlesize': 16, 
                     'xtick.labelsize': 10, 'ytick.labelsize': 10})

fig_rank, ax_rank = plt.subplots(figsize=(14, 8))

colors = plt.cm.Set3(np.linspace(0, 1, len(algorithm_stats)))

bars = ax_rank.bar(range(len(algorithm_stats)), algorithm_stats['Средний ранг'], 
                   color=colors, edgecolor='black', linewidth=1.5)

for i, (bar, val) in enumerate(zip(bars, algorithm_stats['Средний ранг'])):
    ax_rank.text(bar.get_x() + bar.get_width()/2, bar.get_height() + 0.1, 
                 f'{val:.2f}', ha='center', va='bottom', fontsize=10, fontweight='bold')

ax_rank.set_xlabel('Конфигурация', fontsize=14, fontweight='bold')
ax_rank.set_ylabel('Средний ранг', fontsize=14, fontweight='bold')
ax_rank.set_title('Средний ранг конфигураций', fontsize=16, fontweight='bold', pad=20)
ax_rank.set_xticks(range(len(algorithm_stats)))
ax_rank.set_xticklabels(algorithm_stats['Алгоритм'], rotation=45, ha='right')
ax_rank.grid(True, axis='y', linestyle='--', alpha=0.6, linewidth=0.8)
ax_rank.set_axisbelow(True)

plt.tight_layout()
plt.savefig('out/haco/rank.png', dpi=200, bbox_inches='tight', facecolor='white')
plt.show()

rank_pivot = agg.pivot_table(index='name_ru', columns='file_clean', values='rank', aggfunc='first')
order = algorithm_stats.sort_values('Средний ранг')['Алгоритм'].tolist()
rank_pivot = rank_pivot.reindex(order)

fig_heatmap, ax_heatmap = plt.subplots(figsize=(12, 10))

cmap = LinearSegmentedColormap.from_list('RdYlGn', ['green', 'yellow', 'red'], N=100)

im = ax_heatmap.imshow(rank_pivot.values, cmap=cmap, aspect='auto', 
                        interpolation='nearest', vmin=1, vmax=rank_pivot.values.max())

for i in range(rank_pivot.shape[0]):
    for j in range(rank_pivot.shape[1]):
        value = rank_pivot.values[i, j]
        if not pd.isna(value):
            norm_val = (value - 1) / (rank_pivot.values.max() - 1)
            text_color = 'white' if norm_val > 0.5 else 'black'
            ax_heatmap.text(j, i, f'{value:.0f}', ha='center', va='center', 
                          fontsize=10, fontweight='bold', color=text_color)

ax_heatmap.set_xticks(range(len(rank_pivot.columns)))
ax_heatmap.set_xticklabels(rank_pivot.columns, rotation=45, ha='right')
ax_heatmap.set_yticks(range(len(rank_pivot.index)))
ax_heatmap.set_yticklabels(rank_pivot.index, fontsize=9)
ax_heatmap.set_xlabel('Задача TSP', fontsize=14, fontweight='bold')
ax_heatmap.set_ylabel('Конфигурация', fontsize=14, fontweight='bold')
ax_heatmap.set_title('Тепловая карта рангов алгоритмов по задачам', 
                     fontsize=16, fontweight='bold', pad=20)

cbar = plt.colorbar(im, ax=ax_heatmap)
cbar.set_label('Ранг', fontsize=12, fontweight='bold')

plt.tight_layout()
plt.savefig('out/haco/heatmap.png', dpi=200, bbox_inches='tight', facecolor='white')
plt.show()

print("\nГрафики сохранены:")
print("  - out/rank_plot.png/pdf")
print("  - out/rank_heatmap.png/pdf")
import pandas as pd
import matplotlib.pyplot as plt
import matplotlib.patches as mpatches
import numpy as np
import os
from adjustText import adjust_text

# Создаем папку для выходных файлов
os.makedirs('out/convergence', exist_ok=True)

# Читаем данные из файла
df = pd.read_csv('convergence.csv')

# ========== ПЕРЕВОД НАЗВАНИЙ АЛГОРИТМОВ ==========
name_translation = {
    'aco': 'муравьиный',
    'haco_no_local': 'модификация',
    'haco_with_local': 'модификация + лок. оптимизация'
}
df['name_ru'] = df['name'].map(name_translation).fillna(df['name'])

# ========== НАСТРОЙКИ ДЛЯ АЛГОРИТМОВ ==========
algorithms = sorted(df['name_ru'].unique())

# Цвета, штриховки, маркеры для каждого алгоритма
colors = {
    'муравьиный': '#FF6B6B',
    'модификация': '#4ECDC4',
    'модификация + лок. оптимизация': '#90EE90',
}
hatch_patterns = {
    'муравьиный': '/',
    'модификация': '\\',
    'модификация + лок. оптимизация': 'x',
}
markers = {
    'муравьиный': 'o',
    'модификация': 's',
    'модификация + лок. оптимизация': '^',
}

# ========== АГРЕГАЦИЯ ДЛЯ ТАБЛИЦЫ ==========
print("--- Агрегация данных для LaTeX таблицы ---")

# Группируем по file, gensize, name_ru и вычисляем нужные статистики
table_data = df.groupby(['file', 'gensize', 'name_ru']).agg({
    'score': ['mean', 'std'],
    'itb': 'median',
    'auc': 'median',
    'duration_ms': 'mean'
}).round(2)

# Переименовываем колонки
table_data.columns = ['mean_score', 'std_score', 'median_itb', 'median_auc', 'mean_duration']

# Сбрасываем индексы для удобства
table_data = table_data.reset_index()

# ========== СОЗДАНИЕ LaTeX ТАБЛИЦЫ С ИСПОЛЬЗОВАНИЕМ \cline ==========
print("\n--- Создание convergence.tex ---")

latex_lines = []
latex_lines.append("\\begin{longtable}{|p{0.10\\textwidth}|p{0.08\\textwidth}|p{0.14\\textwidth}|p{0.10\\textwidth}|p{0.10\\textwidth}|p{0.10\\textwidth}|p{0.12\\textwidth}|}")
latex_lines.append("\\caption{Результаты сходимости алгоритмов}\\label{tbl:convergence} \\\\")
latex_lines.append("\\hline")
latex_lines.append("\\textbf{Файл} & $N_g$ & \\textbf{Алгоритм} &  $\mu_r$ & $\sigma_r$ & $ITB_{1/2}$ & $\mu_{t}$, мс \\\\")
latex_lines.append("\\hline")
latex_lines.append("\\endfirsthead")
latex_lines.append("")
latex_lines.append("\\caption[]{Результаты сходимости (продолжение)} \\\\")
latex_lines.append("\\hline")
latex_lines.append("\\textbf{Файл} & $N_g$ & \\textbf{Алгоритм} &  $\mu_r$ & $\sigma_r$ & $ITB_{1/2}$ & $\mu_{t}$, мс \\\\")
latex_lines.append("\\hline")
latex_lines.append("\\endhead")
latex_lines.append("\\hline")
latex_lines.append("\\endfoot")
latex_lines.append("\\endlastfoot")

# Группируем по файлу и поколениям
last_file = None
for (file, gensize), group in table_data.groupby(['file', 'gensize']):
    group = group.sort_values('name_ru')
    num_algorithms = len(group)
    file = file.replace('tsp/', '').replace('.tsp', '')

    name_map = {
        'муравьиный': 'мур.',
        'модификация': 'мод.',
        'модификация + лок. оптимизация': 'мод. + лок.',
    }
    
    for idx, (_, row) in enumerate(group.iterrows()):
        name_ru = name_map[row['name_ru']]
        mean_score = row['mean_score']
        std_score = row['std_score']
        median_itb = row['median_itb']
        median_auc = row['median_auc']
        mean_duration = row['mean_duration']
        
        if idx == 0 and last_file != file:
            latex_lines.append(f"\\hline")
            latex_lines.append(f"{file} & {gensize} & {name_ru} & {mean_score:.2f} & {std_score:.2f} & {median_itb:.2f} & {mean_duration:.2f} \\\\")
        elif idx == 0:
            latex_lines.append(f"\\cline{{2-7}}")
            latex_lines.append(f" & {gensize} & {name_ru} & {mean_score:.2f} & {std_score:.2f} & {median_itb:.2f} & {mean_duration:.2f} \\\\")
        else:
            latex_lines.append(f"\\cline{{3-7}}")
            latex_lines.append(f" & & {name_ru} & {mean_score:.2f} & {std_score:.2f} & {median_itb:.2f} & {mean_duration:.2f} \\\\")
    
    last_file = file

latex_lines.append("\\hline")
latex_lines.append("\\end{longtable}")

# Сохраняем в файл
tex_path = 'out/convergence/convergence.tex'
with open(tex_path, 'w', encoding='utf-8') as f:
    f.write('\n'.join(latex_lines))

print(f"✓ LaTeX таблица сохранена в '{tex_path}'")

# Также сохраняем простую CSV версию для проверки
csv_path = 'out/convergence/convergence_table.csv'
table_data.to_csv(csv_path, index=False)
print(f"✓ CSV таблица сохранена в '{csv_path}'")

# ========== АГРЕГАЦИЯ ПО МЕДИАНАМ ==========
print("--- Агрегация данных по медианам ---")
median_data = df.groupby(['file', 'name_ru']).agg({
    'itb': 'median',
    'score': 'median'
}).reset_index()

# ========== ГРАФИК 1: BOXPLOT ITB ==========
print("\n--- Построение boxplot ITB ---")

fig1, ax1 = plt.subplots(figsize=(16, 8))

unique_files = sorted(df['file'].unique())

positions = []
all_data = []
box_colors = []
box_hatches = []
file_centers = {}
file_algo_data = {}
current_pos = 1

for file in unique_files:
    file_positions = []
    for algo in algorithms:
        orig_algo = [k for k, v in name_translation.items() if v == algo][0]
        values = df[(df['file'] == file) & (df['name'] == orig_algo)]['itb'].dropna()
        
        if len(values) > 0:
            all_data.append(values.tolist())
            positions.append(current_pos)
            file_positions.append(current_pos)
            box_colors.append(colors[algo])
            box_hatches.append(hatch_patterns[algo] * 2)
            file_algo_data[(file, algo)] = {'position': current_pos, 'values': values.tolist(), 'color': colors[algo]}
            current_pos += 1
    
    if file_positions:
        file_centers[file] = sum(file_positions) / len(file_positions)
    current_pos += 1

# Строим boxplot
bp = ax1.boxplot(all_data, positions=positions, widths=0.6, patch_artist=True,
                showmeans=True, meanline=True,
                meanprops=dict(linestyle='--', linewidth=2, color='red', alpha=0.8),
                medianprops=dict(linestyle='-', linewidth=2.5, color='black'),
                whiskerprops=dict(color='black', linewidth=1.5),
                capprops=dict(color='black', linewidth=1.5), showfliers=False)

# Применяем цвета и штриховки
for box, color, hatch in zip(bp['boxes'], box_colors, box_hatches):
    box.set_facecolor(color)
    box.set_hatch(hatch)
    box.set_edgecolor('black')
    box.set_linewidth(2)
    box.set_alpha(0.85)

# Добавляем все точки
for (file, algo), data in file_algo_data.items():
    jitter = np.random.normal(data['position'], 0.12, size=len(data['values']))
    ax1.scatter(jitter, data['values'], alpha=0.4, s=25, c=data['color'], 
               edgecolors='black', linewidth=0.5, zorder=3)

# Настройка осей для boxplot
ax1.set_xticks(list(file_centers.values()))
ax1.set_xticklabels(list(file_centers.keys()), fontsize=12, fontweight='bold')
ax1.set_xlabel('Задача', fontsize=14, fontweight='bold')
ax1.set_ylabel('ITB (меньше — лучше)', fontsize=14, fontweight='bold')
ax1.set_title('Распределение ITB по задачам и конфигурациям', fontsize=14, fontweight='bold')
ax1.grid(True, axis='y', linestyle='--', alpha=0.6)
ax1.set_ylim(bottom=0)

# Легенда для boxplot
legend_elements = [mpatches.Patch(facecolor=colors[algo], edgecolor='black', 
                                   hatch=hatch_patterns[algo] * 2, label=algo, alpha=0.85) 
                   for algo in algorithms]
legend_elements.extend([
    plt.Line2D([0], [0], color='red', linestyle='--', linewidth=2, label='Среднее'),
    plt.Line2D([0], [0], color='black', linestyle='-', linewidth=2.5, label='Медиана')
])
ax1.legend(handles=legend_elements, loc='upper right', fontsize=10)

plt.tight_layout()
plt.savefig('out/convergence/itb_boxplots.png', dpi=300, bbox_inches='tight')
print("✓ График 1 сохранен в 'out/convergence/itb_boxplots.png'")
plt.close(fig1)

# ========== ГРАФИК 2: ТОЧЕЧНЫЙ (МЕДИАНА ITB vs МЕДИАНА SCORE) ==========
print("\n--- Построение точечного графика медиан ITB vs Score ---")

fig2, ax2 = plt.subplots(figsize=(12, 9))

texts = []

# Для каждого алгоритма рисуем точки
for algo in algorithms:
    algo_data = median_data[median_data['name_ru'] == algo]
    
    ax2.scatter(
        algo_data['itb'], 
        algo_data['score'],
        alpha=0.7,
        s=120,
        c=colors[algo],
        marker=markers[algo],
        edgecolors='black',
        linewidth=1.5,
        label=algo,
        zorder=3
    )
    
    # Добавляем подписи
    for _, row in algo_data.iterrows():
        text = ax2.annotate(
            row['file'],
            (row['itb'], row['score']),
            fontsize=9,
            alpha=0.8,
            fontweight='bold',
            bbox=dict(boxstyle="round,pad=0.2", facecolor="white", edgecolor="gray", alpha=0.7)
        )
        texts.append(text)

adjust_text(texts, 
    ax=ax2,
    expand_points=(1.5, 1.5),
    expand_text=(1.2, 1.2),
    arrowprops=dict(arrowstyle='->', color='gray', lw=0.5, alpha=0.5),
    force_points=(0.2, 0.2),
    force_text=(0.2, 0.2),
    lim=500)

# Настройка осей
ax2.set_xlabel('ITB (медиана, меньше — лучше)', fontsize=14, fontweight='bold')
ax2.set_ylabel('Длина маршрута (медиана, меньше — лучше)', fontsize=14, fontweight='bold')
ax2.set_title('Зависимость медианы длины маршрута от медианы ITB по файлам', fontsize=14, fontweight='bold')
ax2.grid(True, linestyle='--', alpha=0.6, linewidth=0.8)
ax2.set_xlim(left=0)
ax2.set_ylim(bottom=0)

# Легенда
ax2.legend(loc='upper right', fontsize=12, framealpha=0.95, edgecolor='black')

# Линии для "идеальной" точки
min_itb = median_data['itb'].min()
max_score = median_data['score'].max()
ax2.axvline(x=min_itb, color='gray', linestyle=':', alpha=0.5, linewidth=1)
ax2.axhline(y=max_score, color='gray', linestyle=':', alpha=0.5, linewidth=1)

plt.tight_layout()
plt.savefig('out/convergence/itb_score_scatter.png', dpi=300, bbox_inches='tight')
print("✓ График 2 сохранен в 'out/convergence/itb_score_scatter.png'")
plt.close(fig2)

# ========== ДОПОЛНИТЕЛЬНАЯ ИНФОРМАЦИЯ ==========
print("\n--- Сводная таблица медиан по файлам ---")
pivot_itb = median_data.pivot_table(index='file', columns='name_ru', values='itb')
pivot_score = median_data.pivot_table(index='file', columns='name_ru', values='score')

print("\nМедиана ITB по файлам:")
print(pivot_itb.round(2))

print("\nМедиана Score по файлам:")
print(pivot_score.round(2))

print("\n--- Лучшие алгоритмы по медиане ITB ---")
best_by_file = median_data.loc[median_data.groupby('file')['itb'].idxmin()]
for _, row in best_by_file.iterrows():
    print(f"Файл {row['file']}: {row['name_ru']} (ITB = {row['itb']:.2f})")

print("\n--- Лучшие алгоритмы по медиане Score ---")
best_by_file_score = median_data.loc[median_data.groupby('file')['score'].idxmin()]
for _, row in best_by_file_score.iterrows():
    print(f"Файл {row['file']}: {row['name_ru']} (Score = {row['score']:.2f})")

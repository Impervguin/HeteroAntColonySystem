import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from matplotlib.colors import PowerNorm
import seaborn as sns
from matplotlib.patches import Patch
import os

# Настройка глобальных параметров шрифтов (еще больше увеличены)
plt.rcParams['font.size'] = 14
plt.rcParams['axes.labelsize'] = 16
plt.rcParams['axes.titlesize'] = 16
plt.rcParams['xtick.labelsize'] = 13
plt.rcParams['ytick.labelsize'] = 13
plt.rcParams['legend.fontsize'] = 14
plt.rcParams['figure.titlesize'] = 18

os.makedirs('out/stability', exist_ok=True)

# Чтение данных
df = pd.read_csv('stability.csv')

# ==================== ЧАСТЬ 1: Тепловые карты (две рядом) ====================

def create_combined_heatmap_stability_plot(data_dict, output_path):
    """Создает две тепловые карты рядышком для разных методов"""
    names = list(data_dict.keys())
    fig, axes = plt.subplots(1, 2, figsize=(20, 9))
    
    # Находим общий диапазон для цветовой шкалы
    all_scores = []
    for data in data_dict.values():
        for alpha in sorted(data['alpha'].unique()):
            for beta in sorted(data['beta'].unique()):
                scores = data[(data['alpha'] == alpha) & (data['beta'] == beta)]['score']
                if len(scores) > 0:
                    all_scores.append(scores.mean())
    
    vmin, vmax = min(all_scores), max(all_scores)
    
    # Используем степенную нормализацию (gamma < 1 делает быстрое изменение в начале)
    # gamma=0.5 - квадратный корень, изменение цвета быстрее для малых значений
    # gamma=0.3 - еще более быстрое изменение в начале
    gamma = 0.4  # можно регулировать: чем меньше gamma, тем быстрее изменение в начале
    norm = PowerNorm(gamma=gamma, vmin=vmin, vmax=vmax)
    
    for idx, name in enumerate(names):
        data = data_dict[name]
        ax = axes[idx]
        
        # Получаем уникальные значения
        alphas = sorted(data['alpha'].unique())
        betas = sorted(data['beta'].unique())
        
        # Создаем матрицу средних значений
        mean_matrix = np.zeros((len(betas), len(alphas)))
        std_matrix = np.zeros((len(betas), len(alphas)))
        
        for i, beta in enumerate(betas):
            for j, alpha in enumerate(alphas):
                scores = data[(data['alpha'] == alpha) & (data['beta'] == beta)]['score']
                if len(scores) > 0:
                    mean_matrix[i, j] = scores.mean()
                    std_matrix[i, j] = scores.std()
                else:
                    mean_matrix[i, j] = np.nan
                    std_matrix[i, j] = np.nan
        
        # Рисуем тепловую карту с нелинейной нормализацией
        im = ax.imshow(mean_matrix, cmap='viridis', aspect='auto', 
                       origin='lower', norm=norm)
        
        # Настройка осей
        ax.set_xticks(np.arange(len(alphas)))
        ax.set_yticks(np.arange(len(betas)))
        ax.set_xticklabels([f'{a:.3f}' for a in alphas], rotation=45, ha='right', fontsize=13)
        ax.set_yticklabels([f'{b:.3f}' for b in betas], fontsize=13)
        ax.set_xlabel('α', fontsize=17, fontweight='bold', labelpad=10)
        ax.set_ylabel('β', fontsize=17, fontweight='bold', labelpad=10)
        
        # Добавляем значения в ячейки с инвертированным цветом
        for i in range(len(betas)):
            for j in range(len(alphas)):
                if not np.isnan(mean_matrix[i, j]):
                    # Используем нормализованное значение для определения цвета текста
                    norm_val = norm(mean_matrix[i, j])
                    
                    if norm_val < 0.5:
                        text_color = 'white'  # светлый текст на темном фоне
                    else:
                        text_color = 'black'  # темный текст на светлом фоне
                    
                    ax.text(j, i, f'{mean_matrix[i, j]:.1f}\n±{std_matrix[i, j]:.1f}',
                           ha='center', va='center', fontsize=11, color=text_color, weight='bold')
        
        # Название для подграфика
        name_ru = 'Модифицированный' if name == 'haco' else 'Муравьиный'
        ax.set_title(f'{name_ru} алгоритм\n(среднее ± СКО)', fontsize=17, fontweight='bold', pad=20)
    
    # Общий colorbar - размещаем справа от обоих графиков
    cbar_ax = fig.add_axes([0.92, 0.15, 0.02, 0.7])
    cbar = fig.colorbar(im, cax=cbar_ax)
    cbar.set_label('Средняя длина маршрута', fontsize=15, fontweight='bold')
    cbar.ax.tick_params(labelsize=13)
    
    # Добавляем информацию о нелинейности на colorbar
    cbar.ax.text(1.5, 0.5, f'(γ={gamma})', transform=cbar.ax.transAxes, 
                 fontsize=10, rotation=90, va='center')
    
    # Общий заголовок
    fig.suptitle('Сравнение устойчивости алгоритмов', 
                 fontsize=20, fontweight='bold', y=1.02)
    
    # Настраиваем отступы
    plt.subplots_adjust(left=0.08, right=0.9, wspace=0.25, top=0.88, bottom=0.1)
    
    plt.savefig(output_path, dpi=200, bbox_inches='tight')
    plt.show()
    plt.close()
    print(f"✓ Сохранен совмещенный график (тепловые карты) с нелинейной шкалой: {output_path}")

# Создаем совмещенный график для haco и aco
methods_to_plot = {}
for name in ['haco', 'aco']:
    if name in df['name'].values:
        methods_to_plot[name] = df[df['name'] == name]
    else:
        print(f"⚠ Предупреждение: метод '{name}' не найден в данных")

if methods_to_plot:
    create_combined_heatmap_stability_plot(methods_to_plot, 'out/stability/comparison.png')

# ==================== ЧАСТЬ 2: LaTeX таблица ====================

# Агрегируем данные для всех методов
print("\n" + "="*60)
print("Генерация LaTeX таблицы...")

# Группируем по name, alpha, beta
agg_data = df.groupby(['name', 'alpha', 'beta'])['score'].agg(['mean', 'std']).reset_index()
agg_data['rse'] = (agg_data['std'] / agg_data['mean']) * 100  # RSE в процентах

# Сортируем
agg_data = agg_data.sort_values(['name', 'alpha', 'beta'])

# Генерируем LaTeX код
latex_lines = [
    r"\begin{longtable}{|p{0.12\textwidth}|p{0.10\textwidth}|p{0.12\textwidth}|p{0.12\textwidth}|p{0.12\textwidth}|p{0.12\textwidth}|}",
    r"\caption{Результаты исследования устойчивости алгоритмов}\label{tbl:stability} \\",
    r"\hline",
    r"\textbf{Метод} & $\alpha$ & $\beta$ & \textbf{Среднее} & \textbf{СКО} & \textbf{RSE (\%)} \\",
    r"\hline",
    r"\endfirsthead",
    r"",
    r"\caption*{Продолжение таблицы~\ref{tbl:stability}} \\",
    r"\hline",
    r"\textbf{Метод} & $\alpha$ & $\beta$ & \textbf{Среднее} & \textbf{СКО} & \textbf{RSE (\%)} \\",
    r"\hline",
    r"\endhead",
    r"\hline",
    r"\endfoot",
    r"\endlastfoot",
    r"\hline"
]

# Проходим по данным и формируем строки
current_name = None
last_alpha = None
first_row = True

name_ru_map = {
    'haco': 'Модификация',
    'aco': 'Муравьиный',
}

for idx, row in agg_data.iterrows():
    name = row['name']
    alpha = row['alpha']
    beta = row['beta']
    mean_score = row['mean']
    std = row['std']
    rse = row['rse']
    
    name_ru = name_ru_map[name]
        
    # Формируем строку
    if name != current_name:
        # Первая строка для метода - пишем название
        diff_line = "\\hline"
        line = f"{name_ru} & {alpha:.1f} & {beta:.1f} & {mean_score:.2f} & {std:.2f} & {rse:.2f} \\\\"
        current_name = name
        last_alpha = alpha
    elif alpha != last_alpha:
        diff_line = "\\cline{2-6}"
        line = f" & {alpha:.1f} & {beta:.1f} & {mean_score:.2f} & {std:.2f} & {rse:.2f} \\\\"
        last_alpha = alpha
    else:
        diff_line = "\\cline{3-6}"
        line = f" & & {beta:.1f} & {mean_score:.1f} & {std:.2f} & {rse:.2f} \\\\"
    
    latex_lines.append(diff_line)
    latex_lines.append(line)

latex_lines.append(r"\hline")
latex_lines.append(r"\end{longtable}")

# Сохраняем таблицу
with open('out/stability/results.tex', 'w', encoding='utf-8') as f:
    f.write('\n'.join(latex_lines))

print("✓ LaTeX таблица сохранена в 'out/stability/results.tex'")

# ==================== ЧАСТЬ 3: Размах и CV средних значений ====================

print("\n" + "="*60)
print("Анализ стабильности средних значений:")
print("="*60)

# Для каждого метода сначала считаем средний score для каждой пары alpha/beta
# потом считаем размах и CV этих средних значений

results = []

for name in df['name'].unique():
    name_data = df[df['name'] == name]
    
    # Группируем по alpha и beta, считаем среднее
    mean_scores = name_data.groupby(['alpha', 'beta'])['score'].mean().values
    
    if len(mean_scores) > 1:
        range_val = mean_scores.max() - mean_scores.min()
        cv_val = (mean_scores.std() / mean_scores.mean()) * 100 if mean_scores.mean() != 0 else np.nan
    else:
        range_val = np.nan
        cv_val = np.nan
    
    results.append({
        'method': name,
        'range_mean_scores': range_val,
        'cv_mean_scores': cv_val,
        'n_combinations': len(mean_scores)
    })

# Выводим результаты
print("\nМетод | Размах средних | CV средних (%) | Кол-во комбинаций")
print("-" * 70)

for res in results:
    print(f"{res['method']:10} | {res['range_mean_scores']:14.2f} | {res['cv_mean_scores']:14.2f} | {res['n_combinations']}")

# Сохраняем результаты в CSV
results_df = pd.DataFrame(results)
results_df.to_csv('out/stability/mean_scores_stability.csv', index=False)
print("\n✓ Результаты сохранены в 'out/stability/mean_scores_stability.csv'")

# Детальный вывод для отчета
print("\n" + "="*60)
print("ДЕТАЛЬНЫЙ АНАЛИЗ:")
print("="*60)

for res in results:
    name_ru = name_ru_map.get(res['method'], res['method'].upper())
    print(f"\n📊 {name_ru} алгоритм:")
    print(f"   • Размах средних значений: {res['range_mean_scores']:.2f}")
    print(f"   • Коэффициент вариации (CV): {res['cv_mean_scores']:.2f}%")
    
    # Интерпретация CV
    if res['cv_mean_scores'] < 10:
        print(f"   • Стабильность: Низкая вариабельность (стабильный метод)")
    elif res['cv_mean_scores'] < 25:
        print(f"   • Стабильность: Средняя вариабельность")
    else:
        print(f"   • Стабильность: Высокая вариабельность (нестабильный метод)")

print("\n" + "="*60)
print("✅ Все задачи выполнены!")
print("  - Совмещенный график (тепловые карты): out/stability/comparison.png")
print("  - LaTeX таблица: out/stability/results.tex")
print("  - CSV с анализом: out/stability/mean_scores_stability.csv")
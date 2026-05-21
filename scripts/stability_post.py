import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D
from matplotlib import cm
from matplotlib.patches import Patch
import os

os.makedirs('out/stability', exist_ok=True)

# Чтение данных
df = pd.read_csv('stability.csv')

# ==================== ЧАСТЬ 1: 3D графики (два в одном) ====================

def create_combined_3d_stability_plot(data_dict, output_path):
    """Создает два 3D столбчатых графика рядышком для разных методов"""
    names = list(data_dict.keys())
    fig = plt.figure(figsize=(16, 8))
    z_mean_max = 0
    z_mean_min = float('inf')
    for idx, name in enumerate(names, 1):
        data = data_dict[name]
        alphas = sorted(data['alpha'].unique())
        betas = sorted(data['beta'].unique())
        for alpha in alphas:
            for beta in betas:
                scores = data[(data['alpha'] == alpha) & (data['beta'] == beta)]['score']
                z_mean = scores.mean()
                if z_mean > z_mean_max:
                    z_mean_max = z_mean
                if z_mean < z_mean_min:
                    z_mean_min = z_mean
    
    for idx, name in enumerate(names, 1):
        data = data_dict[name]
        ax = fig.add_subplot(1, 2, idx, projection='3d')

        # Получаем уникальные значения
        alphas = sorted(data['alpha'].unique())
        betas = sorted(data['beta'].unique())
        
        # Создаем сетку для столбцов
        x_pos = np.arange(len(alphas))
        y_pos = np.arange(len(betas))
        X, Y = np.meshgrid(x_pos, y_pos)
        X = X.flatten()
        Y = Y.flatten()
        
        # Собираем данные
        z_mean = []
        z_std = []
        
        for alpha in alphas:
            for beta in betas:
                scores = data[(data['alpha'] == alpha) & (data['beta'] == beta)]['score']
                if len(scores) > 0:
                    z_mean.append(scores.mean())
                    z_std.append(scores.std())
                else:
                    z_mean.append(0)
                    z_std.append(0)
        
        Z = np.array(z_mean)
        errors = np.array(z_std)
        
        # Цветовая карта на основе значений
        norm = plt.Normalize(z_mean_min, z_mean_max)
        colors = cm.viridis(norm(Z))
        
        # Рисуем столбцы
        bars = ax.bar3d(X, Y, np.zeros_like(Z), dx=0.7, dy=0.7, dz=Z, 
                        color=colors, alpha=0.7, edgecolor='black', linewidth=0.5)
        
        # Добавляем планки ошибок (std)
        for i, (x, y, z, err) in enumerate(zip(X, Y, Z, errors)):
            # Верхняя планка
            ax.plot([x, x], [y, y], [z, z + err], color='red', linewidth=2, alpha=0.8)
            ax.plot([x-0.15, x+0.15], [y, y], [z + err, z + err], color='red', linewidth=2, alpha=0.8)
            # Нижняя планка
            if z - err >= 0:
                ax.plot([x, x], [y, y], [z, z - err], color='red', linewidth=2, alpha=0.8)
                ax.plot([x-0.15, x+0.15], [y, y], [z - err, z - err], color='red', linewidth=2, alpha=0.8)
        
        # Настройка осей
        ax.set_xticks(x_pos)
        ax.set_xticklabels([f'{a:.3f}' for a in alphas], rotation=45, ha='right', fontsize=8)
        ax.set_yticks(y_pos)
        ax.set_yticklabels([f'{b:.3f}' for b in betas], fontsize=12)
        ax.set_xlabel('α', fontsize=14, labelpad=10)
        ax.set_ylabel('β', fontsize=14, labelpad=10)
        ax.set_zlabel('Средняя длина маршрута', fontsize=14, labelpad=10)
        z_min = ax.get_zlim()[0]
        ax.set_zlim([z_min, z_mean_max])
        
        # Название для подграфика
        name_ru = 'Модификация' if name == 'haco' else 'Муравьиный'
        ax.set_title(f'{name_ru} алгоритм', fontsize=13, fontweight='bold', pad=20)
        
        # Настройка угла обзора для лучшей видимости
        ax.view_init(elev=25, azim=-60)
    
    # Общий заголовок
    fig.suptitle('Сравнение устойчивости алгоритмов: средняя длина маршрута', 
                 fontsize=14, fontweight='bold', y=0.98)
    
    # Добавляем общий colorbar
    # Собираем все значения для общего colorbar
    all_scores = []
    for data in data_dict.values():
        for alpha in sorted(data['alpha'].unique()):
            for beta in sorted(data['beta'].unique()):
                scores = data[(data['alpha'] == alpha) & (data['beta'] == beta)]['score']
                if len(scores) > 0:
                    all_scores.append(scores.mean())
    
    if all_scores:
        sm = plt.cm.ScalarMappable(cmap=cm.viridis, 
                                   norm=plt.Normalize(min(all_scores), max(all_scores)))
        sm.set_array([])
        cbar = fig.colorbar(sm, ax=fig.axes, shrink=0.5, aspect=20, pad=0.05)
        cbar.set_label('Средняя длина маршрута', fontsize=12)
    
    plt.tight_layout()
    plt.savefig(output_path, dpi=200, bbox_inches='tight')
    plt.show()
    plt.close()
    print(f"✓ Сохранен совмещенный график: {output_path}")

# Создаем совмещенный график для haco и aco
methods_to_plot = {}
for name in ['haco', 'aco']:
    if name in df['name'].values:
        methods_to_plot[name] = df[df['name'] == name]
    else:
        print(f"⚠ Предупреждение: метод '{name}' не найден в данных")

if methods_to_plot:
    create_combined_3d_stability_plot(methods_to_plot, 'out/stability/comparison.png')

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
print("  - Совмещенный график: out/stability/comparison.png")
print("  - LaTeX таблица: out/stability/results.tex")
print("  - CSV с анализом: out/stability/mean_scores_stability.csv")
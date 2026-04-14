import pandas as pd
import plotly.express as px

# === 1. Загрузка ===
df = pd.read_csv("results.csv")

# === 2. Агрегация по (алгоритм, файл) ===
agg = df.groupby(["algorithm", "file"]).agg({
    "duration_ms": ["mean", "std"],
    "score": ["mean", "std"]
}).reset_index()

# Упрощаем названия колонок
agg.columns = [
    "algorithm", "file",
    "time_mean", "time_std",
    "score_mean", "score_std"
]

# === 3. График времени ===
fig_time = px.bar(
    agg,
    x="file",
    y="time_mean",
    color="algorithm",
    error_y="time_std",
    barmode="group",
    title="Execution Time (ms)",
    labels={
        "file": "TSP Instance",
        "time_mean": "Time (ms)",
        "algorithm": "Algorithm"
    }
)

fig_time.update_layout(xaxis_tickangle=-30)
fig_time.show()

# === 4. График качества ===
fig_score = px.bar(
    agg,
    x="file",
    y="score_mean",
    color="algorithm",
    error_y="score_std",
    barmode="group",
    title="Solution Quality (lower is better)",
    labels={
        "file": "TSP Instance",
        "score_mean": "Score",
        "algorithm": "Algorithm"
    }
)

fig_score.update_layout(xaxis_tickangle=-30)
fig_score.show()

# === 5. (опционально) сохранить ===
fig_time.write_html("time_plot.html")
fig_score.write_html("score_plot.html")
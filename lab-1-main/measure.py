import os
import time
import subprocess
import matplotlib.pyplot as plt
import numpy as np

plt.rcParams.update({
    "font.size": 10,
    "axes.titlesize": 14,
    "axes.labelsize": 12,
    "legend.fontsize": 10,
    "figure.titlesize": 16
})

sizes_mb = []
sizes_bytes = []
times = []

# 生成数据：输入文件大小（以 MB 为单位），记录运行时间
for i in range(1, 21):
    input_file = f"data/input{i}.dat"
    output_file = f"data/output{i}.dat"

    if not os.path.exists(input_file):
        print(f"File not found: {input_file}")
        continue

    size_bytes = os.path.getsize(input_file)
    size_mb = round(size_bytes / (1024 * 1024), 2)  # 转换成 MB，保留两位小数

    print(f"Sorting input{i}.dat ({size_mb} MB)...")

    start = time.perf_counter()
    subprocess.run(["bin\\sort.exe", input_file, output_file], check=True)
    end = time.perf_counter()

    elapsed = round(end - start, 4) 
    sizes_bytes.append(size_bytes)
    sizes_mb.append(size_mb)
    times.append(elapsed)

# 转 numpy 方便计算
sizes_np = np.array(sizes_bytes)
times_np = np.array(times)

# O(n log n) 参考线（缩放到实际运行时间的范围）
complexity = sizes_np * np.log2(sizes_np)
complexity_norm = complexity / np.max(complexity) * np.max(times_np)

# 理论曲线：O(n^2)
complexity_n2 = sizes_np ** 2
complexity_n2 = complexity_n2 / np.max(complexity_n2) * np.max(times_np)

# 绘图
plt.figure(figsize=(10, 6))
plt.plot(sizes_mb, times_np, label="Actual Runtime", marker='o', color='#007acc', linewidth=2)
plt.plot(sizes_mb, complexity_norm, label="O(n log n)", linestyle='--', color='red', linewidth=2)
plt.plot(sizes_mb, complexity_n2, label="O(n²)", linestyle=':', color='green', linewidth=2)

# # 添加每个点的坐标注释
# for x, y in zip(sizes_mb, times_np):
#     plt.text(x, y, f"({x}, {y})", fontsize=8, ha='center', color='black')

plt.xlabel("Input Size (MB)")
plt.ylabel("Time (seconds)")
plt.title("Sort Program's Runtime vs Input Size")
plt.legend()
plt.grid(True)
plt.ylim(0, max(times_np) * 1.2)
plt.tight_layout()
plt.savefig("report.pdf")
print("✅ report.pdf generated with labeled points.")

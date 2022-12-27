
global m_acc
global m_inf_time
global m_size
global m_down_time
global m_addr

m1_acc = 0.05147
m1_inf_time = 2.862
m1_size = 1155 #KB
m1_down_time = 9.88
m1_addr = "https://drive.google.com/drive/folders/1bYnFDHWmanEr0BZOrzAnqMCpv7iDKqLM?usp=sharing"

m2_acc = 0.04867
m2_inf_time = 2.867
m2_size = 1621
m2_down_time = 9.64
m2_addr = "https://drive.google.com/drive/folders/1WEQ0IH4iu7NxCQ9PZ-G1o8AMuKkppe5B?usp=sharing"

m3_acc = 0.04161
m3_inf_time = 2.895
m3_size = 2181
m3_down_time = 13.45
m3_addr = "https://drive.google.com/drive/folders/1zdpyXbnptSU4Yb0KdtJOrAUoXb3DN3EU?usp=sharing"

m4_acc = 0.03849
m4_inf_time = 2.873
m4_size = 2815
m4_down_time = 10.43
m4_addr = "https://drive.google.com/drive/folders/1JwJ31dFTEHIWzlwL3gTe8idtxYY0tOET?usp=sharing"

m5_acc = 0.03392
m5_inf_time = 2.8790
m5_size = 4336
m5_down_time = 2.52
m5_addr = "https://drive.google.com/drive/folders/1SRktV1_rYbfVGVPXjBu0KiFLDq8i7uhJ?usp=sharing"


m_acc = [m1_acc, m2_acc, m3_acc, m4_acc, m5_acc]
m_inf_time = [m1_inf_time, m2_inf_time, m3_inf_time, m4_inf_time, m5_inf_time]
m_size = [m1_size, m2_size, m3_size, m4_size, m5_size]
m_down_time = [m1_down_time, m2_down_time, m3_down_time, m4_down_time, m5_down_time]
m_completion_time = []
m_completion_time = [m_down_time[i] + m_inf_time[i] for i in range(len(m_down_time))]
m_addr = [m1_addr, m2_addr, m3_addr, m4_addr, m5_addr]
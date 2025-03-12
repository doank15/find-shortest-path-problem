# Cấu trúc dữ liệu:
- Đồ thị được biểu diễn bằng danh sách kề (adj), trong đó mỗi phần tử chứa danh sách các cạnh từ một đỉnh.
Mỗi cạnh (Edge) chứa thông tin về đỉnh đích và trọng số.
Lớp Graph hỗ trợ cả đồ thị có hướng và vô hướng.

# Thuật toán tìm đường đi ngắn nhất:
- `Dijkstra`: Sử dụng cho đồ thị không có cạnh âm. Thuật toán này dùng hàng đợi ưu tiên (priority queue) để chọn đỉnh có khoảng cách nhỏ nhất trong mỗi bước.
- `SPFA (Shortest Path Faster Algorithm)`: Dùng cho đồ thị có cạnh âm nhưng không có chu trình âm. Thuật toán này sử dụng hàng đợi thông thường và có khả năng phát hiện chu trình âm.
- `Johnson`: Giải quyết bài toán tìm đường đi ngắn nhất giữa mọi cặp đỉnh trong đồ thị. Thuật toán này biến đổi trọng số cạnh để loại bỏ cạnh âm, sau đó áp dụng Dijkstra từ mỗi đỉnh.

# Chiến lược chọn thuật toán:
- Hàm `FindShortestPath` tự động kiểm tra đồ thị có cạnh âm hay không.
- Nếu đồ thị không có cạnh âm, sử dụng `Dijkstra` để đạt hiệu suất tốt nhất.
- Nếu đồ thị có cạnh âm, sử dụng `SPFA` để xử lý đúng và phát hiện chu trình âm.
# Xử lý trường hợp đặc biệt:
- Phát hiện chu trình âm trong đồ thị.
- Xử lý trường hợp không tồn tại đường đi.
- Kiểm tra tính hợp lệ của các đỉnh đầu vào.
# Giải pháp này hiệu quả vì:
- Tự động chọn thuật toán tối ưu dựa trên đặc điểm của đồ thị.
- Xử lý được nhiều tình huống khác nhau: đồ thị có/không có cạnh âm, có/không có chu trình âm.
- Sử dụng các cấu trúc dữ liệu hiệu quả như hàng đợi ưu tiên cho `Dijkstra`.
- Đảm bảo tính đúng đắn bằng cách phát hiện và xử lý các trường hợp đặc biệt.
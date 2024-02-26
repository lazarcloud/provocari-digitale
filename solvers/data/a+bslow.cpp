#include <iostream>
#include <numeric>  // Include the <numeric> header for std::accumulate
#include <vector>
int A[100000000];
int main() {
  // Allocate a large vector
  for (int i = 0; i < 100000000; i++) {
    A[i] = 7;
  }
  int x, y;
  std::cin >> x >> y;
  std::cout << x + y;
  return 0;
}

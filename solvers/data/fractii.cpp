#include <iostream>

// Funcție pentru calculul CMMDC (gcd)
int gcd(int a, int b) {
  while (b) {
    int r = a % b;
    a = b;
    b = r;
  }
  return a;
}

// Funcție pentru calculul numărului de fracții ireductibile
int countIrreducibleFractions(int N) {
  int count = 0;
  for (int i = 1; i <= N; ++i) {
    for (int j = 1; j <= N; ++j) {
      if (gcd(i, j) == 1) {
        ++count;
      }
    }
  }
  return count;
}

int main() {
  int N;
  std::cin >> N;

  int result = countIrreducibleFractions(N);
  std::cout << result;

  return 0;
}

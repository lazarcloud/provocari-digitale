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

int main() {
  int num1, num2;
  std::cin >> num1 >> num2;

  // Calculăm CMMDC
  int cmmdc = gcd(num1, num2);

  // Verificăm dacă cele două numere sunt prime între ele
  if (cmmdc == 1) {
    std::cout << 0;
  } else {
    std::cout << cmmdc;
  }

  return 0;
}

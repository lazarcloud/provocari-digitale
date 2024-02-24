#include <fstream>
#include <iostream>
#include <numeric>
#include <vector>
using namespace std;
int A[100000000];
int main() {
  ifstream fin("input.in");
  ofstream fout("output.out");
  // Allocate a large vector
  for (int i = 0; i < 100000000; i++) {
    A[i] = 7;
  }
  int x, y;
  fin >> x >> y;
  fout << x + y;
  return 0;
}

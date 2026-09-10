package matrix_test

import (
	"math"
	"testing"

	"go-api/internal/matrix"
)

const tolerance = 1e-9

func multiply(A, B [][]float64) [][]float64 {
	m := len(A)
	n := len(A[0])
	p := len(B[0])

	C := make([][]float64, m)
	for i := 0; i < m; i++ {
		C[i] = make([]float64, p)
		for j := 0; j < p; j++ {
			sum := 0.0
			for k := 0; k < n; k++ {
				sum += A[i][k] * B[k][j]
			}
			C[i][j] = sum
		}
	}
	return C
}

func transpose(A [][]float64) [][]float64 {
	m := len(A)
	n := len(A[0])

	At := make([][]float64, n)
	for j := 0; j < n; j++ {
		At[j] = make([]float64, m)
		for i := 0; i < m; i++ {
			At[j][i] = A[i][j]
		}
	}
	return At
}

func TestGramSchmidtQR_SquareMatrix(t *testing.T) {
	A := [][]float64{
		{1.0, 2.0, 4.0},
		{3.0, 8.0, 14.0},
		{2.0, 6.0, 13.0},
	}

	Q, R, err := matrix.GramSchmidtQR(A)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m, n := len(A), len(A[0])

	// 1. Verify Q shape (m x n) and R shape (n x n)
	if len(Q) != m || len(Q[0]) != n {
		t.Errorf("Q shape mismatch: got %dx%d, want %dx%d", len(Q), len(Q[0]), m, n)
	}
	if len(R) != n || len(R[0]) != n {
		t.Errorf("R shape mismatch: got %dx%d, want %dx%d", len(R), len(R[0]), n, n)
	}

	// 2. Verify Q * R == A
	QR := multiply(Q, R)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if math.Abs(QR[i][j]-A[i][j]) > tolerance {
				t.Errorf("Q*R != A at [%d][%d]: got %f, want %f", i, j, QR[i][j], A[i][j])
			}
		}
	}

	// 3. Verify Q^T * Q == I_n
	Qt := transpose(Q)
	QtQ := multiply(Qt, Q)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			expected := 0.0
			if i == j {
				expected = 1.0
			}
			if math.Abs(QtQ[i][j]-expected) > tolerance {
				t.Errorf("Q^T*Q != I at [%d][%d]: got %f, want %f", i, j, QtQ[i][j], expected)
			}
		}
	}

	// 4. Verify R is upper triangular (R[i][j] == 0 for i > j)
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			if math.Abs(R[i][j]) > tolerance {
				t.Errorf("R is not upper triangular at [%d][%d]: got %f", i, j, R[i][j])
			}
		}
	}
}

func TestGramSchmidtQR_RectangularTall(t *testing.T) {
	// 3x2 matrix (m > n)
	A := [][]float64{
		{1.0, 2.0},
		{3.0, 4.0},
		{5.0, 6.0},
	}

	Q, R, err := matrix.GramSchmidtQR(A)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	QR := multiply(Q, R)
	for i := 0; i < len(A); i++ {
		for j := 0; j < len(A[0]); j++ {
			if math.Abs(QR[i][j]-A[i][j]) > tolerance {
				t.Errorf("Q*R != A at [%d][%d]: got %f, want %f", i, j, QR[i][j], A[i][j])
			}
		}
	}
}

func TestGramSchmidtQR_Errors(t *testing.T) {
	_, _, err := matrix.GramSchmidtQR([][]float64{})
	if err == nil {
		t.Error("expected error for empty matrix")
	}

	_, _, err = matrix.GramSchmidtQR([][]float64{{1.0, 2.0}, {3.0}})
	if err == nil {
		t.Error("expected error for non-rectangular matrix")
	}
}

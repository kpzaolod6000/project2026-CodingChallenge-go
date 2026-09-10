package matrix

import (
	"errors"
	"math"
)

var (
	ErrEmptyMatrix    = errors.New("la matriz no puede estar vacía")
	ErrNonRectangular = errors.New("la matriz debe ser rectangular")
)

// GramSchmidtQR calcula la descomposición QR de una matriz A (m x n) utilizando
// el algoritmo de Gram-Schmidt Modificado (MGS).
// Devuelve la matriz ortogonal Q (m x n) y la matriz triangular superior R (n x n)
// tal que A = Q * R.
func GramSchmidtQR(A [][]float64) ([][]float64, [][]float64, error) {
	if len(A) == 0 || len(A[0]) == 0 {
		return nil, nil, ErrEmptyMatrix
	}

	m := len(A)
	n := len(A[0])

	for _, row := range A {
		if len(row) != n {
			return nil, nil, ErrNonRectangular
		}
	}

	// Extraer vectores columna de A
	v := make([][]float64, n)
	for j := 0; j < n; j++ {
		v[j] = make([]float64, m)
		for i := 0; i < m; i++ {
			v[j][i] = A[i][j]
		}
	}

	// Inicializar columnas de Q y la matriz R
	qCols := make([][]float64, n)
	for j := 0; j < n; j++ {
		qCols[j] = make([]float64, m)
	}

	R := make([][]float64, n)
	for i := 0; i < n; i++ {
		R[i] = make([]float64, n)
	}

	// Algoritmo de Gram-Schmidt Modificado
	for k := 0; k < n; k++ {
		u := make([]float64, m)
		copy(u, v[k])

		for i := 0; i < k; i++ {
			dot := 0.0
			for r := 0; r < m; r++ {
				dot += qCols[i][r] * u[r]
			}
			R[i][k] = dot

			for r := 0; r < m; r++ {
				u[r] -= dot * qCols[i][r]
			}
		}

		norm := 0.0
		for r := 0; r < m; r++ {
			norm += u[r] * u[r]
		}
		norm = math.Sqrt(norm)

		R[k][k] = norm

		if norm > 1e-12 {
			for r := 0; r < m; r++ {
				qCols[k][r] = u[r] / norm
			}
		} else {
			for r := 0; r < m; r++ {
				qCols[k][r] = 0.0
			}
		}
	}

	// Reconstruir Q en formato orientado a filas (m x n)
	Q := make([][]float64, m)
	for i := 0; i < m; i++ {
		Q[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			Q[i][j] = qCols[j][i]
		}
	}

	return Q, R, nil
}

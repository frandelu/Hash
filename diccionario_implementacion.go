package diccionario

import (
	"fmt"
	"hash/fnv"
	//Link de la funcion de hash https://pkg.go.dev/hash/fnv?utm_source=gopls
)

const (
	_VACIO           = 0
	_OCUPADO         = 1
	_BORRADO         = -1
	_CAP_INCIAL      = 32
	_FACTOR_AGRANDAR = 2
	_FACTOR_ACHICAR  = 4
	_PORCENTAJE      = 100
	_CARGA_MAX       = 65
	_CARGA_MIN       = 20
)

type elemento[K comparable, V any] struct {
	estado int
	clave  K
	dato   V
}
type diccionario[K comparable, V any] struct {
	elementos []elemento[K, V]
	capacidad int
	cantidad  int
	borrados  int
}

type iterDiccionario[K comparable, V any] struct {
	dicc       *diccionario[K, V]
	pos_actual int
}

func CrearHash[K comparable, V any]() Diccionario[K, V] {
	dicc := new(diccionario[K, V])
	dicc.capacidad = _CAP_INCIAL
	dicc.elementos = make([]elemento[K, V], _CAP_INCIAL)
	return dicc
}

// Primitivas Diccionario

func (dicc *diccionario[K, V]) Guardar(clave K, dato V) {
	dicc.redimension()
	pos := dicc.calcularPos(clave)
	if dicc.elementos[pos].estado == _VACIO {
		dicc.cantidad++
	}
	dicc.elementos[pos].estado = _OCUPADO
	dicc.elementos[pos].clave = clave
	dicc.elementos[pos].dato = dato
}

func (dicc *diccionario[K, V]) Pertenece(clave K) bool {
	pos := dicc.calcularPos(clave)
	return dicc.elementos[pos].estado == _OCUPADO
}

func (dicc *diccionario[K, V]) Obtener(clave K) V {
	if !dicc.Pertenece(clave) {
		panic("La clave no pertenece al diccionario")
	}
	pos := dicc.calcularPos(clave)
	return dicc.elementos[pos].dato
}

func (dicc *diccionario[K, V]) Borrar(clave K) V {
	dicc.redimension()
	if !dicc.Pertenece(clave) {
		panic("La clave no pertenece al diccionario")
	}
	pos := dicc.calcularPos(clave)
	dicc.elementos[pos].estado = _BORRADO
	dicc.cantidad--
	dicc.borrados++
	return dicc.elementos[pos].dato
}

func (dicc *diccionario[K, V]) Cantidad() int {
	return dicc.cantidad
}

func (dicc *diccionario[K, V]) Iterar(visitar func(clave K, dato V) bool) {
	for i := 0; i < dicc.capacidad; i++ {
		elem := dicc.elementos[i]
		if elem.estado == _OCUPADO && !visitar(elem.clave, elem.dato) {
			return
		}
	}
}

func (dicc *diccionario[K, V]) Iterador() IterDiccionario[K, V] {
	iterador := new(iterDiccionario[K, V])
	iterador.dicc = dicc
	iterador.pos_actual = 0
	if dicc.elementos[0].estado != _OCUPADO {
		iterador.pos_actual = iterador.buscarSiguiente()
	}
	return iterador
}

//Primitivas de IterDiccionario

func (iter *iterDiccionario[K, V]) HaySiguiente() bool {
	for i := iter.pos_actual; i < iter.dicc.capacidad; i++ {
		if iter.dicc.elementos[i].estado == _OCUPADO {
			return true
		}
	}
	return false
}

func (iter *iterDiccionario[K, V]) VerActual() (K, V) {
	if !iter.HaySiguiente() {
		panic("El iterador termino de iterar")
	}
	return iter.dicc.elementos[iter.pos_actual].clave, iter.dicc.elementos[iter.pos_actual].dato
}

func (iter *iterDiccionario[K, V]) Siguiente() K {
	if !iter.HaySiguiente() {
		panic("El iterador termino de iterar")
	}
	clave_act := iter.dicc.elementos[iter.pos_actual].clave
	iter.pos_actual = iter.buscarSiguiente()
	return clave_act
}

// Funciones / métodos auxiliares

func hash(clave []byte) uint64 {
	x := fnv.New64a()
	x.Write(clave)
	return (x.Sum64)()
}

func convertirABytes[K comparable](clave K) []byte {
	return []byte(fmt.Sprintf("%v", clave))
}

func (dicc *diccionario[K, V]) calcularPos(clave K) uint64 {
	pos := hash(convertirABytes(clave)) % uint64(dicc.capacidad)
	for dicc.elementos[pos].estado != _VACIO {
		if dicc.elementos[pos].clave == clave && dicc.elementos[pos].estado == _OCUPADO {
			return pos
		}
		if pos+1 == uint64(dicc.capacidad) {
			pos = 0
		} else {
			pos++
		}
	}
	return pos
}

func (dicc *diccionario[K, V]) redimension() {
	carga := ((dicc.cantidad + dicc.borrados) * _PORCENTAJE) / dicc.capacidad
	vieja_cap := dicc.capacidad
	if carga > _CARGA_MAX {
		dicc.capacidad *= _FACTOR_AGRANDAR
	} else if carga < _CARGA_MIN {
		dicc.capacidad /= _FACTOR_ACHICAR
		if _CAP_INCIAL > dicc.capacidad {
			dicc.capacidad = _CAP_INCIAL
		}
	} else {
		return
	}
	elementos := dicc.elementos
	dicc.elementos = make([]elemento[K, V], dicc.capacidad)
	dicc.borrados = 0

	for i := 0; i < vieja_cap; i++ {
		if elementos[i].estado != _OCUPADO {
			continue
		}
		pos := dicc.calcularPos(elementos[i].clave)
		dicc.elementos[pos] = elementos[i]
	}
}

func (iter *iterDiccionario[K, V]) buscarSiguiente() int {
	for i := iter.pos_actual + 1; i < iter.dicc.capacidad; i++ {
		if iter.dicc.elementos[i].estado == _OCUPADO {
			return i
		}
	}
	return iter.dicc.capacidad
}

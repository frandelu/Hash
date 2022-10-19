package diccionario

import (
	"fmt"
	"hash/fnv"
	//Link de la funcion de hash https://pkg.go.dev/hash/fnv?utm_source=gopls
)

const (
	_VACIO = iota
	_OCUPADO
	_BORRADO

	_CAP_INCIAL   = 32
	_FACTOR_REDIM = 2
	_CARGA_MAX    = 75
	_CARGA_MIN    = 20
)

type elemento[K comparable, V any] struct {
	estado int
	clave  K
	dato   V
}
type hashCerrado[K comparable, V any] struct {
	elementos []elemento[K, V]
	capacidad int
	cantidad  int
	borrados  int
}

type iterHashCerrado[K comparable, V any] struct {
	dicc       *hashCerrado[K, V]
	pos_actual int
}

func CrearHash[K comparable, V any]() Diccionario[K, V] {
	dicc := new(hashCerrado[K, V])
	dicc.capacidad = _CAP_INCIAL
	dicc.elementos = make([]elemento[K, V], _CAP_INCIAL)
	return dicc
}

// Primitivas Diccionario

func (dicc *hashCerrado[K, V]) Guardar(clave K, dato V) {
	carga := ((dicc.cantidad + dicc.borrados) * 100) / dicc.capacidad
	if carga > _CARGA_MAX {
		dicc.redimension(dicc.capacidad * _FACTOR_REDIM)
	}
	pos := dicc.calcularPos(clave)
	if dicc.elementos[pos].estado == _VACIO {
		dicc.cantidad++
	}
	dicc.elementos[pos].estado = _OCUPADO
	dicc.elementos[pos].clave = clave
	dicc.elementos[pos].dato = dato
}

func (dicc *hashCerrado[K, V]) Pertenece(clave K) bool {
	pos := dicc.calcularPos(clave)
	return dicc.elementos[pos].estado == _OCUPADO
}

func (dicc *hashCerrado[K, V]) Obtener(clave K) V {
	pos := dicc.calcularPos(clave)
	if dicc.elementos[pos].estado == _OCUPADO {
		return dicc.elementos[pos].dato
	} else {
		panic("La clave no pertenece al diccionario")
	}
}

func (dicc *hashCerrado[K, V]) Borrar(clave K) V {
	carga := (dicc.cantidad * 100) / dicc.capacidad
	if carga < _CARGA_MIN && dicc.capacidad > _CAP_INCIAL {
		dicc.redimension(dicc.capacidad / _FACTOR_REDIM)
	}
	pos := dicc.calcularPos(clave)
	if dicc.elementos[pos].estado == _OCUPADO {
		dicc.elementos[pos].estado = _BORRADO
		dicc.cantidad--
		dicc.borrados++
	} else {
		panic("La clave no pertenece al diccionario")
	}
	return dicc.elementos[pos].dato
}

func (dicc *hashCerrado[K, V]) Cantidad() int {
	return dicc.cantidad
}

func (dicc *hashCerrado[K, V]) Iterar(visitar func(clave K, dato V) bool) {
	for i := 0; i < dicc.capacidad; i++ {
		elem := dicc.elementos[i]
		if elem.estado == _OCUPADO && !visitar(elem.clave, elem.dato) {
			return
		}
	}
}

func (dicc *hashCerrado[K, V]) Iterador() IterDiccionario[K, V] {
	iterador := new(iterHashCerrado[K, V])
	iterador.dicc = dicc
	iterador.pos_actual = 0
	if dicc.elementos[0].estado != _OCUPADO {
		iterador.pos_actual = iterador.buscarSiguiente()
	}
	return iterador
}

//Primitivas de IterDiccionario

func (iter *iterHashCerrado[K, V]) HaySiguiente() bool {
	/*for i := iter.pos_actual; i < iter.dicc.capacidad; i++ {
		if iter.dicc.elementos[i].estado == _OCUPADO {
			return true
		}
	}
	return false*/
	return iter.pos_actual != iter.dicc.capacidad
}

func (iter *iterHashCerrado[K, V]) VerActual() (K, V) {
	if !iter.HaySiguiente() {
		panic("El iterador termino de iterar")
	}
	return iter.dicc.elementos[iter.pos_actual].clave, iter.dicc.elementos[iter.pos_actual].dato
}

func (iter *iterHashCerrado[K, V]) Siguiente() K {
	if iter.pos_actual == iter.dicc.capacidad {
		panic("El iterador termino de iterar")
	}
	clave_act := iter.dicc.elementos[iter.pos_actual].clave
	iter.pos_actual = iter.buscarSiguiente()
	return clave_act
}

// Funciones / métodos auxiliares

func convertirABytes[K comparable](clave K) []byte {
	return []byte(fmt.Sprintf("%v", clave))
}

func hash(clave []byte) uint64 {
	x := fnv.New64a()
	x.Write(clave)
	return x.Sum64()
}

func (dicc *hashCerrado[K, V]) calcularPos(clave K) uint64 {
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

func (dicc *hashCerrado[K, V]) redimension(nueva_cap int) {
	vieja_cap := dicc.capacidad
	elementos := dicc.elementos
	dicc.elementos = make([]elemento[K, V], nueva_cap)
	dicc.borrados = 0
	dicc.capacidad = nueva_cap
	for i := 0; i < vieja_cap; i++ {
		if elementos[i].estado != _OCUPADO {
			continue
		}
		pos := dicc.calcularPos(elementos[i].clave)
		dicc.elementos[pos] = elementos[i]
	}
}

func (iter *iterHashCerrado[K, V]) buscarSiguiente() int {
	for i := iter.pos_actual + 1; i < iter.dicc.capacidad; i++ {
		if iter.dicc.elementos[i].estado == _OCUPADO {
			return i
		}
	}
	return iter.dicc.capacidad
}

package diccionario

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
	//Link de la funcion de hash https://pkg.go.dev/hash/fnv?utm_source=gopls
)

const (
	_VACIO           = 0
	_OCUPADO         = 1
	_BORRADO         = -1
	_CAP_INCIAL      = 8
	_FACTOR_AGRANDAR = 2
	_FACTOR_ACHICAR  = 4
	_PORCETAJE       = 100
	_CARGA_MAX       = 50
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
	diccionario diccionario[K, V]
}

func hash(clave []byte) uint64 {
	x := fnv.New64a()
	x.Write(clave)
	return (x.Sum64)()
}

func convertirABytes[K comparable](clave K) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(clave)
	return buf.Bytes()
}

func calcularPos[K comparable](clave K, capacidad int) uint64 {
	return (hash(convertirABytes(clave))) % uint64(capacidad)
}

func redimension[K comparable, V any](dicc *diccionario[K, V]) {
	carga := ((dicc.cantidad + dicc.borrados) * _PORCETAJE) / dicc.capacidad
	if carga > _CARGA_MAX {
		var nuevo_dicc diccionario[K, V]
		nuevo_dicc.capacidad = dicc.capacidad * _FACTOR_AGRANDAR
		nuevo_dicc.elementos = make([]elemento[K, V], nuevo_dicc.capacidad)
		nuevo_dicc.borrados = dicc.borrados
		nuevo_dicc.cantidad = dicc.cantidad
		for i := 0; i < dicc.capacidad; i++ {
			if dicc.elementos[i].estado == _OCUPADO {
				nuevo_dicc.Guardar(dicc.elementos[i].clave, dicc.elementos[i].dato)
			}
		}
	}
	if carga > _CARGA_MIN {
		var nuevo_dicc diccionario[K, V]
		nuevo_dicc.capacidad = dicc.capacidad / _FACTOR_ACHICAR
		nuevo_dicc.elementos = make([]elemento[K, V], nuevo_dicc.capacidad)
		nuevo_dicc.borrados = dicc.borrados
		nuevo_dicc.cantidad = dicc.cantidad
		for i := 0; i < dicc.capacidad; i++ {
			if dicc.elementos[i].estado == _OCUPADO {
				nuevo_dicc.Guardar(dicc.elementos[i].clave, dicc.elementos[i].dato)
			}
		}
	}

}

func CrearHash[K comparable, V any]() diccionario[K, V] {
	dicc := new(diccionario[K, V])
	dicc.borrados = 0
	dicc.cantidad = 0
	dicc.capacidad = _CAP_INCIAL
	dicc.elementos = make([]elemento[K, V], _CAP_INCIAL)
	return *dicc
}

func (dicc *diccionario[K, V]) Guardar(clave K, dato V) {
	redimension(dicc)
	pos := calcularPos(clave, dicc.capacidad)
	var pos_final uint64
	for pos_final = pos; dicc.elementos[pos_final].estado == _BORRADO ||
		(dicc.elementos[pos_final].clave != clave && dicc.elementos[pos_final].estado == _OCUPADO); pos_final++ {
		if pos_final == uint64(dicc.capacidad)-1 {
			pos_final = 0
		}
	}
	dicc.elementos[pos_final].clave = clave
	dicc.elementos[pos_final].dato = dato
	dicc.elementos[pos_final].estado = _OCUPADO
	dicc.cantidad++
}

func (dicc *diccionario[K, V]) Pertenece(clave K) bool {
	pos := calcularPos(clave, dicc.capacidad)
	for i := 0; dicc.elementos[pos].estado != _VACIO && i <= dicc.cantidad; i++ {
		if dicc.elementos[pos].clave == clave {
			return true
		}
		if pos == uint64(dicc.capacidad)-1 {
			pos = 0
			continue
		}
		pos++
	}
	return false
}

func (dicc *diccionario[K, V]) Obtener(clave K) V {
	if !dicc.Pertenece(clave) {
		panic("La clave no pertenece al diccionario")
	}
	pos := calcularPos(clave, dicc.capacidad)
	for dicc.elementos[pos].estado != _VACIO {
		if dicc.elementos[pos].clave == clave {
			break
		}
		if pos == uint64(dicc.capacidad)-1 {
			pos = 0
			continue
		}
		pos++
	}
	return dicc.elementos[pos].dato
}

func (dicc *diccionario[K, V]) Borrar(clave K) V {
	if !dicc.Pertenece(clave) {
		panic("La clave no pertenece al diccionario")
	}
	pos := calcularPos(clave, dicc.capacidad)
	for dicc.elementos[pos].estado != _VACIO {
		if dicc.elementos[pos].clave == clave {
			dicc.elementos[pos].estado = _BORRADO
			dicc.cantidad--
			dicc.borrados++
			break
		}
		if pos == uint64(dicc.capacidad)-1 {
			pos = 0
			continue
		}
		pos++
	}
	return dicc.elementos[pos].dato
}

func (dicc *diccionario[K, V]) Cantidad() int {
	return dicc.cantidad
}

func (dicc *diccionario[K, V]) Iterar(visitar func(clave K, dato V) bool) {
	for i := 0; visitar(dicc.elementos[i].clave, dicc.elementos[i].dato) && i < dicc.capacidad; i++ {
	}
}

//Primitivas del Iterador externo

func (dicc *diccionario[K, V]) Iterador() iterDiccionario[K, V] {
	iterador := new(iterDiccionario[K, V])
	return *iterador
}

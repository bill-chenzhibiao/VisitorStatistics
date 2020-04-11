package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFloat64Obj2IntSuccess(t *testing.T){
	var f float64 = 3
	var obj interface{} = f
	assert.Equal(t,3,obj2int(obj,0))
}


func TestObj2IntFailForDefault(t *testing.T){
	var f int = 3
	var obj interface{} = f
	assert.Equal(t,0,obj2int(obj,0))
}
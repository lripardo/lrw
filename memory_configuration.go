package lrw

import (
	"encoding/json"
	"log"
	"strconv"
)

type MemoryConfiguration struct {
	Data map[string]interface{}
}

func (m *MemoryConfiguration) String(key ConfigurationKey) string {
	data := m.Data[key.Key]
	if data == nil {
		return key.Default
	}
	value, ok := data.(string)
	if !ok {
		return key.Default
	}
	return value
}

func (m *MemoryConfiguration) Uint(key ConfigurationKey) uint {
	data := m.Data[key.Key]
	if data == nil {
		n, err := strconv.ParseUint(key.Default, 10, 32)
		if err != nil {
			log.Panic(err)
		}
		return uint(n)
	}
	value, ok := data.(uint)
	if !ok {
		n, err := strconv.ParseUint(key.Default, 10, 32)
		if err != nil {
			log.Panic(err)
		}
		return uint(n)
	}
	return value
}

func (m *MemoryConfiguration) Bool(key ConfigurationKey) bool {
	data := m.Data[key.Key]
	if data == nil {
		b, err := strconv.ParseBool(key.Default)
		if err != nil {
			log.Panic(err)
		}
		return b
	}
	value, ok := data.(bool)
	if !ok {
		b, err := strconv.ParseBool(key.Default)
		if err != nil {
			log.Panic(err)
		}
		return b
	}
	return value
}

func (m *MemoryConfiguration) Int(key ConfigurationKey) int {
	data := m.Data[key.Key]
	if data == nil {
		n, err := strconv.ParseInt(key.Default, 10, 32)
		if err != nil {
			log.Panic(err)
		}
		return int(n)
	}
	value, ok := data.(int)
	if !ok {
		n, err := strconv.ParseInt(key.Default, 10, 32)
		if err != nil {
			log.Panic(err)
		}
		return int(n)
	}
	return value
}

func (m *MemoryConfiguration) Int64(key ConfigurationKey) int64 {
	data := m.Data[key.Key]
	if data == nil {
		n, err := strconv.ParseInt(key.Default, 10, 64)
		if err != nil {
			log.Panic(err)
		}
		return n
	}
	value, ok := data.(int64)
	if !ok {
		n, err := strconv.ParseInt(key.Default, 10, 64)
		if err != nil {
			log.Panic(err)
		}
		return n
	}
	return value
}

func (m *MemoryConfiguration) Strings(key ConfigurationKey) []string {
	data := m.Data[key.Key]
	if data == nil {
		var array []string
		if err := json.Unmarshal([]byte(key.Default), &array); err != nil {
			log.Panic(err)
		}
		return array
	}
	value, ok := data.([]string)
	if !ok {
		var array []string
		if err := json.Unmarshal([]byte(key.Default), &array); err != nil {
			log.Panic(err)
		}
		return array
	}
	return value
}

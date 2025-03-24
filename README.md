# Go API Units Example

<p align="center">
  <em>Handling measurement units in Go APIs - the right way</em>
</p>

<p align="center">
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat-square&logo=go" alt="Go Version"></a>
  <a href="https://github.com/haimkastner/go-api-units-example/blob/main/LICENSE"><img src="https://img.shields.io/github/license/haimkastner/go-api-units-example?style=flat-square" alt="License"></a>
  <a href="https://units-docs.gleece.dev/"><img src="https://img.shields.io/badge/API-Documentation-ff69b4?style=flat-square" alt="API Docs"></a>
</p>

## 📖 Overview

This repository demonstrates best practices for representing and working with measurement units in Go APIs. Gone are the days of "is this temperature in Celsius or Fahrenheit?" confusion in your API!

Using this approach, your API can:
- ✅ Accept input values in any compatible unit
- ✅ Internally work with standardized units
- ✅ Return values in just units
- ✅ Provide clear, self-documenting interfaces

## 🛠️ Key Dependencies

- [unitsnet-go](https://github.com/haimkastner/unitsnet-go) - Comprehensive unit conversion library with 100+ unit types
- [Gleece](https://github.com/gleece/gleece) - Go framework for building APIs from code.

## 🚀 Live Demo

Check out the [OpenAPI 3.0 Documentation](https://units-docs.gleece.dev/) live demonstration

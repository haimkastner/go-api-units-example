package controllers

import (
	"github.com/gopher-fleece/runtime"
	"github.com/haimkastner/unitsnet-go/units"
)

// UnitsController
// @Tag(Units) Units Operations
// @Route(/units)
// @Description The Units API Example
type UnitsController struct {
	runtime.GleeceController // Embedding the GleeceController to inherit its methods
}

// The LengthFactory object to create the Length objects
var lf = units.LengthFactory{}

// @Description Post unit API and return the processed unit
// @Method(POST)
// @Route(/post-unit)
// @Query(useUnit) The unit to be used in response
// @Body(data) The unit to process
// @Response(200) The response with the processed unit
// @ErrorResponse(500) The error when process failed
func (ec *UnitsController) TestUnit(useUnit *units.LengthUnits, data units.LengthDto) (units.LengthDto, error) {
	// The unit to be processed
	var unit *units.Length

	// Load the unit from the DTO
	unit, _ = units.LengthFactory{}.FromDto(data)

	// TODO: Process the unit (logic here)

	// Return the processed unit
	return unit.ToDto(useUnit), nil
}

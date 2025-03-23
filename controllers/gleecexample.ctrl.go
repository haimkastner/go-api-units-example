package controllers

import (
	"github.com/gopher-fleece/runtime"
	"github.com/haimkastner/unitsnet-go/units"
)

// UsersController
// @Tag(Units) Users
// @Route(/units)
// @Description The Units API Example
type UnitsController struct {
	runtime.GleeceController // Embedding the GleeceController to inherit its methods
}

// @Description Post unit API and return the processed unit
// @Method(POST)
// @Route(/post-unit)
// @Query(useUnit) The unit to be used in response
// @Body(data) The unit to process
// @Response(200) The response with the processed unit
// @ErrorResponse(500) The error when process failed
func (ec *UnitsController) TestUnit(useUnit *units.LengthUnits, data units.LengthDto) (units.LengthDto, error) {
	lf := units.LengthFactory{}
	unit, _ := lf.FromDto(data)
	// return error if the unit is not valid

	// Do the logic here using the abstracted unit

	return unit.ToDto(useUnit), nil
}

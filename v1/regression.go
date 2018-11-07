//Package regression provides online linear regression calculation.
package regression

import (
	"math"
)

//Regression represents a queue of past points. Use New() to initialize.
type Regression struct {
	xSum, ySum, xxSum, xySum, yySum, xDelta            float64
	lastSlopeCalc, lastInterceptCalc, lastStdErrorCalc float64
	N                                                  int

	//here so multiple calcs calls per add calls wont hurt performance
	lastCalcFresh bool
}

//New returns a Regression that keeps points back as far as xDelta from the last
//added point.
func New() Regression {
	return Regression{}
}

//Calculate returns the slope, intercept and standard error of a best fit line to the added
//points. Returns a cached value if called between adds. Deprecated in favor of CalculateWithStdError.
func (r *Regression) Calculate() (slope, intercept float64) {
	slope, intercept, _ = r.CalculateWithStdError()
	return
}

//CalculateWithStdError returns the slope, intercept and standard error of a best fit line to the added
//points. Returns a cached value if called between adds.
func (r *Regression) CalculateWithStdError() (slope, intercept, stdError float64) {
	if r.lastCalcFresh {
		slope = r.lastSlopeCalc
		intercept = r.lastInterceptCalc
		stdError = r.lastStdErrorCalc
		return
	}

	n := float64(r.N)

	//linear regression formula:
	//slope is (n*SUM(x*y) - SUM(x)*SUM(y)) / (n*SUM(x*x) - (SUM(x))^2)
	//intercept is (SUM(y)-slope*SUM(x)) / n
	xSumOverN := r.xSum / n //here to only calc once for performance
	slope = (r.xySum - xSumOverN*r.ySum) / (r.xxSum - xSumOverN*r.xSum)
	intercept = (r.ySum - slope*r.xSum) / n

	//standard error formula is sqrt(SUM((yActual - yPredicted)^2) / (n - 2))
	//the n-2 is related to the degrees of freedom for the regression, 2 in this case
	//simplification of the sum
	//SUM((yA - yP)^2)
	//SUM(yA*yA - 2*yA*yP + yP*yP)
	//SUM(y*y) - SUM(2*y*(m*x+b)) + SUM((m*x+b)(m*x+b))
	//SUM(y*y) - 2*m*SUM(x*y) - 2*b*SUM(y) + m*m*SUM(x*x) + 2*b*m*SUM(x) + n*b*b
	twoTimesB := 2 * intercept
	stdError = math.Sqrt((r.yySum - 2*slope*r.xySum - twoTimesB*r.ySum + slope*slope*r.xxSum + twoTimesB*slope*r.xSum + n*intercept*intercept) / (n - 2))

	r.lastSlopeCalc = slope
	r.lastInterceptCalc = intercept
	r.lastStdErrorCalc = stdError
	r.lastCalcFresh = true
	return
}

//Add adds the new x and y as a point into the queue. Panics if given an x value less than the last.
func (r *Regression) Add(x, y float64) {
	r.lastCalcFresh = false
	r.N++
	r.xSum += x
	r.ySum += y
	r.xxSum += x * x
	r.xySum += x * y
	r.yySum += y * y
}

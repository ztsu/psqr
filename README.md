### psqr

### WIP

The P^2 Algorithm for Dynamic Calculation of Quantiles and Histograms Without Storing Observations
(http://www.cs.wustl.edu/~jain/papers/ftp/psqr.pdf)

**How many data you can store on your machine?**

**And how many additional space (RAM/disk) you need for quantile estimation?**

The main problem is that all N observations must be stored.
In many situations, N can be very large; also, there may be many variables whose quantiles may be required.

**It is this space problem.**

Instead of storing the complete distribution function,

**we store only five points on it and update the points as more observations are generated.**

Today and now written on golang!

### Installation:
```bash
go get github.com/theairkit/psqr
```

### Description and examples:

It's deadly simple:
```go
// Define type, satisfying psqr.Sort interface:
type Value float64
func (value Value) Value()float64{
return float64(value)
}
// Set percentile and create new psqr.PSQR:
percentile:=0.95
somePSQR:=psqr.New(percentile)
// Add next value:
nextVal:=Value(12.34)
psqr.AddValue(nextVal)
// Get current quantile at any moment:
quantile:=somePSQR.Quantile()
```

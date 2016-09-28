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


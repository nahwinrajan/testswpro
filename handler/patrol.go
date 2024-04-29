package handler

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
)

// Ideally this is put in Config files, but since this is test
// and we'll move things around better to play it safe and make it
// as mere const for time being
const (
	distanceBetweenPlan = 10
	treeHeightMax       = 30
	treeHeightMin       = 1
	monitorDistance     = 1

	// stepStrFormat: step_#,x,y,direction,step_distance,current_distance
	stepStrFormat = "%d,%d,%d,%s,%d,%d"
)

const (
	directionEW = "ew" // to the left x axis
	directionWE = "we" // to the right x axis
	directionSN = "sn" // up on y axis
	directionNS = "ns" // down on y axis
	directionVU = "vu" // up adjusting for tree height
	directionVD = "vd" // down adjusting for tree height
)

// NOTE: ideally I would want to put this in "usecase" layer
// but since this is following SDK and there is no such layer,
// let's just follow as is

// patrol do patrol the estate and calculate values for the stats and distance
func (srv *Server) patrol(ctx context.Context, estateID string) error {
	estate, err := srv.Repository.GetEstateByID(
		ctx,
		estateID,
	)
	if err != nil {
		return err
	}

	trees, err := srv.Repository.GetAllTreesInEstate(ctx, estateID)
	if err != nil {
		return err
	}

	// we are making it so it is base 1 instead index 0
	fields := make([][]int, 0, estate.Length+1)
	for colIdx := 0; colIdx < estate.Length+1; colIdx++ {
		fields[colIdx] = make([]int, 0, estate.Width+1)
	}

	// TODO: in case performance detoriate, do this within the distance / route loop
	// populate the fields with respective trees in their plot
	lenTrees := len(trees)
	heights := make([]int, 0, lenTrees)
	var minHeight, maxHeight = math.MaxInt, math.MinInt
	for _, tree := range trees {
		fields[tree.Y][tree.X] = tree.Height

		heights = append(heights, tree.Height)
		if tree.Height > maxHeight {
			maxHeight = tree.Height
		} else if tree.Height < minHeight {
			minHeight = tree.Height
		}
	}

	var medianHeight int
	sort.Ints(heights)
	if lenTrees%2 == 0 {
		medianHeight = (heights[lenTrees/2-1] + heights[lenTrees/2]) / 2
	} else {
		medianHeight = heights[lenTrees/2]
	}

	var currDistance, verticalMove int
	var currDirection string
	var strb strings.Builder
	currHeight, currStepNumber := 1, 1

	// stepStrFormat: step_#,x,y,direction,step_distance,current_distance
	// stepStrFormat = "%d,%d,%d,%s,%d,%d"

	for y := 1; y <= estate.Length; y++ {
		switch {
		case y%2 == 0:
			// if its even row, loop from column end to column start
			for x := estate.Width; x > 0; x-- {
				// calculate vertical movement distance if there is tree planted
				// as discussed adjust to tree height is higher priority to avoid drone crashing
				if fields[y][x] != 0 {
					switch {
					case fields[y][x] < currHeight:
						verticalMove = (currHeight - fields[y][x]) + monitorDistance
						currDirection = directionVD
					case fields[y][x] > currHeight:
						verticalMove = (fields[y][x] - currHeight) + monitorDistance
						currDirection = directionVU
					}

					currDistance = currDistance + verticalMove
					currHeight = fields[y][x] + monitorDistance

					fmt.Fprintf(&strb, stepStrFormat, currStepNumber, x, y, currDirection, verticalMove, currDistance)
					currStepNumber++
				}

				currDistance += distanceBetweenPlan
				switch {
				case y != 1 && x == 1:
					currDirection = directionSN
				default:
					currDirection = directionEW
				}

				fmt.Fprintf(&strb, stepStrFormat, currStepNumber, x, y, currDirection, distanceBetweenPlan, currDistance)
				currStepNumber++
			}
		default:
			// if its odd row, loop from column start to column end
			for x := 1; x <= estate.Width; x++ {
				// calculate vertical movement distance if there is tree planted
				// as discussed adjust to tree height is higher priority to avoid drone crashing
				if fields[y][x] != 0 {
					switch {
					case fields[y][x] < currHeight:
						verticalMove = (currHeight - fields[y][x]) + monitorDistance
						currDirection = directionVD
					case fields[y][x] > currHeight:
						verticalMove = (fields[y][x] - currHeight) + monitorDistance
						currDirection = directionVU
					}

					currDistance = currDistance + verticalMove
					currHeight = fields[y][x] + monitorDistance

					fmt.Fprintf(&strb, stepStrFormat, currStepNumber, x, y, currDirection, verticalMove, currDistance)
					currStepNumber++
				}

				currDistance += distanceBetweenPlan
				switch {
				case y != 1 && x == 1:
					currDirection = directionSN
				default:
					currDirection = directionWE
				}

				fmt.Fprintf(&strb, stepStrFormat, currStepNumber, x, y, currDirection, distanceBetweenPlan, currDistance)
				currStepNumber++
			}
		}
	}

	err = srv.Repository.UpdateEstate(
		ctx,
		lenTrees,
		minHeight,
		maxHeight,
		medianHeight,
		currDistance,
		strb.String(),
	)

	return err
}

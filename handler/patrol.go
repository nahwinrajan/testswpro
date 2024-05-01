package handler

import (
	"context"
	"fmt"
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

	// stepStrFormat: step_#,x,y,direction,step_distance,current_distance
	stepStrFormat = "%d,%d,%d,%s,%d,%d;"
)

const (
	directionEW         = "ew" // to the left x axis
	directionWE         = "we" // to the right x axis
	directionSN         = "sn" // up on y axis
	directionNS         = "ns" // down on y axis
	directionVU         = "vu" // up adjusting for tree height
	directionVD         = "vd" // down adjusting for tree height
	directionSameHeight = "--" // down adjusting for tree height
)

// NOTE: ideally I would want to put this in "usecase" layer
// but since this is following SDK and there is no such layer,
// let's just follow as is

// patrol do patrol the estate and calculate values for the stats and distance
func (srv *Server) patrol(ctx context.Context, estateID string) error {
	estate, err := srv.repository.GetEstateByID(
		ctx,
		estateID,
	)
	if err != nil {
		return err
	}

	trees, err := srv.repository.GetAllTreesInEstate(ctx, estateID)
	if err != nil {
		return err
	}

	lenTrees := len(trees)
	if lenTrees < 1 {
		return nil
	}

	// we are making it so it is base 1 instead index 0
	lenRow := estate.Length + 1
	lenCol := estate.Width + 1
	fields := make([][]int, lenRow)
	for rowIdx := 0; rowIdx < lenRow; rowIdx++ {
		fields[rowIdx] = make([]int, lenCol)
	}

	// TODO: in case performance detoriate, do this within the distance / route loop
	// populate the fields with respective trees in their plot
	heights := make([]int, 0, lenTrees)
	var minHeight, maxHeight int
	for _, tree := range trees {
		fields[tree.Y][tree.X] = tree.Height

		heights = append(heights, tree.Height)
		// if this is the first iteration
		if minHeight == 0 && maxHeight == 0 {
			maxHeight = tree.Height
			minHeight = tree.Height
			continue
		}

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
	currDroneHeight, currTreeHeight, currStepNumber := 0, 0, 0
	var monitorDistance = 1

	// stepStrFormat: step_#,x,y,direction,step_distance,current_distance
	// stepStrFormat = "%d,%d,%d,%s,%d,%d;"

	for y := 1; y < lenRow; y++ {
		// calculate the distance between row on row change
		if y > 1 {
			currDistance += distanceBetweenPlan
			currDirection = directionSN // always go up on changing row
			currStepNumber++

			fmt.Fprintf(&strb, stepStrFormat, currStepNumber, estate.Width, y, currDirection, distanceBetweenPlan, currDistance)
		}

		// crude hack to mimick drone movement
		if y%2 == 0 {
			// if it is even row, loop from end (width) to start / west to east
			// the tree given by input is 1 index base not 0 index base, thus we just leave the 0 row and col empty
			for x := lenCol - 1; x > 0; x++ {
				currTreeHeight = fields[y][x]
				// calculate vertical movement distance if there is tree planted
				if currTreeHeight > 0 {
					if currDroneHeight == (currTreeHeight + monitorDistance) {
						verticalMove = 0
						currDirection = directionSameHeight
					} else if (currTreeHeight + monitorDistance) < currDroneHeight {
						verticalMove = (currDroneHeight - (currTreeHeight + monitorDistance))
						currDirection = directionVD
					} else if (currTreeHeight + monitorDistance) > currDroneHeight {
						verticalMove = ((currTreeHeight + monitorDistance) - currDroneHeight)
						currDirection = directionVU
					}

					currDistance = currDistance + verticalMove
					currDroneHeight = currTreeHeight + monitorDistance
					currStepNumber++

					fmt.Fprintf(&strb, stepStrFormat, currStepNumber, x, y, currDirection, verticalMove, currDistance)
				}

				currDistance += distanceBetweenPlan
				currDirection = directionWE
				currStepNumber++
				fmt.Fprintf(&strb, stepStrFormat, currStepNumber, x, y, currDirection, verticalMove, currDistance)
			}
		} else {
			// if its odd row, loop from column start to column end (east to west)
			for x := 1; x < lenCol; x++ {
				currTreeHeight = fields[y][x]
				// calculate vertical movement distance if there is tree planted
				if currTreeHeight > 0 {
					if currDroneHeight == (currTreeHeight + monitorDistance) {
						verticalMove = 0
						currDirection = directionSameHeight
					} else if (currTreeHeight + monitorDistance) < currDroneHeight {
						verticalMove = (currDroneHeight - (currTreeHeight + monitorDistance))
						currDirection = directionVD
					} else if (currTreeHeight + monitorDistance) > currDroneHeight {
						verticalMove = ((currTreeHeight + monitorDistance) - currDroneHeight)
						currDirection = directionVU
					}

					currDistance = currDistance + verticalMove
					currDroneHeight = currTreeHeight + monitorDistance
					currStepNumber++

					fmt.Fprintf(&strb, stepStrFormat, currStepNumber, x, y, currDirection, verticalMove, currDistance)
				}

				currDistance += distanceBetweenPlan
				currDirection = directionEW
				currStepNumber++
				fmt.Fprintf(&strb, stepStrFormat, currStepNumber, x, y, currDirection, verticalMove, currDistance)
			}
		}
	}

	if currDroneHeight > 0 {
		currDistance = currDistance + currDroneHeight
	}

	err = srv.repository.UpdateEstate(
		ctx,
		estateID,
		lenTrees,
		minHeight,
		maxHeight,
		medianHeight,
		currDistance,
		strb.String(),
	)

	return err
}

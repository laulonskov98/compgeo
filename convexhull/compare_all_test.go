package convexhull

import (
	"fmt"
	"testing"
)

func TestSquarePoints(t *testing.T) {
	points := Generate_square_inputs(10400)

	hull := INC_CH(points)
	hull_parallel := PAR_GS(points, 6)
	hull_giftwrap := GIFT_CH(points)

	fmt.Println("giftwrap:")
	fmt.Println(hull_giftwrap)
	fmt.Println("parallel:")
	fmt.Println(hull_parallel)
	fmt.Println("normal:")
	fmt.Println(hull)

	if len(hull) == 0 || len(hull_parallel) == 0 || len(hull_giftwrap) == 0 {
		t.Errorf("One of the hulls is empty, got the following length (%d, %d, %d)", len(hull), len(hull_parallel), len(hull_giftwrap))
	}

	if len(hull) != len(hull_parallel) || len(hull) != len(hull_giftwrap) || len(hull_parallel) != len(hull_giftwrap) {
		t.Errorf("The number of points in the convex hulls are different, got the following length (%d, %d, %d)", len(hull), len(hull_parallel), len(hull_giftwrap))
	}

	for i := 0; i < len(hull); i++ {
		if hull[i] != hull_parallel[i] || hull[i] != hull_giftwrap[i] || hull_parallel[i] != hull_giftwrap[i] {
			t.Errorf("The points in the convex hulls are different, got the following points (%v, %v, %v)", hull[i], hull_parallel[i], hull_giftwrap[i])
		}
	}

}

func TestCirclePoints(t *testing.T) {
	points := Generate_circle_inputs(100)

	hull := INC_CH(points)
	hull_parallel := PAR_GS(points, 6)
	hull_giftwrap := GIFT_CH(points)

	if len(hull) == 0 || len(hull_parallel) == 0 || len(hull_giftwrap) == 0 {
		t.Errorf("One of the hulls is empty, got the following length (%d, %d, %d)", len(hull), len(hull_parallel), len(hull_giftwrap))
	}

	if len(hull) != len(hull_parallel) || len(hull) != len(hull_giftwrap) || len(hull_parallel) != len(hull_giftwrap) {
		t.Errorf("The number of points in the convex hulls are different, got the following length (%d, %d, %d)", len(hull), len(hull_parallel), len(hull_giftwrap))
	}

	for i := 0; i < len(hull); i++ {
		if hull[i] != hull_parallel[i] || hull[i] != hull_giftwrap[i] || hull_parallel[i] != hull_giftwrap[i] {
			t.Errorf("The points in the convex hulls are different, got the following points (%v, %v, %v)", hull[i], hull_parallel[i], hull_giftwrap[i])
		}
	}

}

func TestPolynomialPoints(t *testing.T) {
	points := Generate_polynomial_inputs(20000)

	hull := INC_CH(points)
	hull_parallel := PAR_GS(points, 6)
	hull_giftwrap := GIFT_CH(points)

	fmt.Println(len(hull), len(hull_parallel), len(hull_giftwrap))
	if len(hull) == 0 || len(hull_parallel) == 0 || len(hull_giftwrap) == 0 {
		t.Errorf("One of the hulls is empty, got the following length (%d, %d, %d)", len(hull), len(hull_parallel), len(hull_giftwrap))
	}

	if len(hull) != len(hull_parallel) || len(hull) != len(hull_giftwrap) || len(hull_parallel) != len(hull_giftwrap) {
		t.Errorf("The number of points in the convex hulls are different, got the following length (%d, %d, %d)", len(hull), len(hull_parallel), len(hull_giftwrap))
	}

	for i := 0; i < len(hull); i++ {
		if hull[i] != hull_parallel[i] || hull[i] != hull_giftwrap[i] || hull_parallel[i] != hull_giftwrap[i] {
			t.Errorf("The points in the convex hulls are different, got the following points (%v, %v, %v)", hull[i], hull_parallel[i], hull_giftwrap[i])
		}
	}

}

func TestIncreasingPolynomialSides(t *testing.T) {

	for i := 5; i < 20; i++ {
		points := Generate_n_side_polygon(i, 2000)
		fmt.Println(points)

		hull := INC_CH(points)
		hull_parallel := PAR_GS(points, 4)
		hull_giftwrap := GIFT_CH(points)

		fmt.Println(len(hull), len(hull_parallel), len(hull_giftwrap))
		if len(hull) == 0 || len(hull_parallel) == 0 || len(hull_giftwrap) == 0 {
			t.Errorf("One of the hulls is empty, got the following length (%d, %d, %d)", len(hull), len(hull_parallel), len(hull_giftwrap))
		}

		if len(hull) != len(hull_parallel) || len(hull) != len(hull_giftwrap) || len(hull_parallel) != len(hull_giftwrap) {
			t.Errorf("The number of points in the convex hulls are different, got the following length (%d, %d, %d)", len(hull), len(hull_parallel), len(hull_giftwrap))
		}

		for i := 0; i < len(hull); i++ {
			if hull[i] != hull_parallel[i] || hull[i] != hull_giftwrap[i] || hull_parallel[i] != hull_giftwrap[i] {
				t.Errorf("The points in the convex hulls are different, got the following points (%v, %v, %v)", hull[i], hull_parallel[i], hull_giftwrap[i])
			}
		}
	}
}

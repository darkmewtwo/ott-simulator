package cognition

import (
	"math"
	"reflect"
	"testing"

	"simulator/internal/user"
)

func TestApplyRate(t *testing.T) {
	for _, tc := range []struct {
		name              string
		value, rate, want float64
	}{
		{"increase toward one", .5, .2, .6},
		{"decrease proportionally", .5, -.2, .4},
		{"zero rate", .5, 0, .5},
		{"clamp initial low", -2, .2, .2},
		{"clamp initial high", 2, -.2, .8},
		{"clamp result low", .5, -3, 0},
		{"clamp result high", .5, 3, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := applyRate(tc.value, tc.rate); math.Abs(got-tc.want) > 1e-12 {
				t.Fatalf("applyRate = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestApplyMentalStateCost(t *testing.T) {
	state := user.MentalState{Energy: .5, Happiness: .5, Boredom: .5, Frustration: .5, Immersion: .5, Satisfaction: .5}
	ApplyMentalStateCost(&state, MentalStateCost{Energy: -.2, Happiness: .2, Boredom: -.4, Frustration: .4, Immersion: .6, Satisfaction: -.6})
	want := user.MentalState{Energy: .4, Happiness: .6, Boredom: .3, Frustration: .7, Immersion: .8, Satisfaction: .2}
	if !reflect.DeepEqual(state, want) {
		t.Fatalf("state = %+v, want %+v", state, want)
	}
}

func TestCostForAction(t *testing.T) {
	for _, tc := range []struct {
		action CatalogAction
		want   MentalStateCost
	}{
		{ActionSelectMovie, MentalStateCost{Energy: -.01, Happiness: .05, Boredom: -.20, Frustration: -.10, Satisfaction: .10}},
		{ActionContinueBrowsing, MentalStateCost{Energy: -.04, Happiness: -.01, Boredom: .06, Frustration: .04, Satisfaction: -.03}},
		{ActionLeaveSession, MentalStateCost{}},
		{"unknown", MentalStateCost{}},
	} {
		if got := CostForAction(tc.action); got != tc.want {
			t.Errorf("cost for %s = %+v, want %+v", tc.action, got, tc.want)
		}
	}
}

func TestCalculateCatalogScores(t *testing.T) {
	ctx := CatalogContext{
		Personality: user.Personality{Impulsiveness: 1, Curiosity: 1, Exploration: 1, Patience: 1, AttentionSpan: 1},
		MentalState: user.MentalState{Energy: 1, Happiness: 1, Satisfaction: 1},
	}
	scores := CalculateCatalogScores(ctx)
	want := []ActionScore{{ActionSelectMovie, .8}, {ActionContinueBrowsing, 1}, {ActionLeaveSession, -1}}
	if len(scores) != len(want) {
		t.Fatalf("got %d scores", len(scores))
	}
	for i, score := range scores {
		if score.Action != want[i].Action || math.Abs(score.Score-want[i].Score) > 1e-12 {
			t.Errorf("score = %+v, want %+v", score, want[i])
		}
	}
	t.Run("custom scorers receive context in order", func(t *testing.T) {
		calls := 0
		scorer := func(got CatalogContext) ActionScore {
			if got != ctx {
				t.Error("context was not forwarded")
			}
			calls++
			return ActionScore{ActionLeaveSession, float64(calls)}
		}
		got := CalculateCatalogScores(ctx, scorer, scorer)
		if !reflect.DeepEqual(got, []ActionScore{{ActionLeaveSession, 1}, {ActionLeaveSession, 2}}) {
			t.Fatalf("custom scores = %+v", got)
		}
	})
}

func TestFrustratedExhaustedUserPrefersLeaving(t *testing.T) {
	ctx := CatalogContext{MentalState: user.MentalState{Frustration: 1, Boredom: 1}}
	leave := ScoreLeaveSession(ctx)
	if math.Abs(leave.Score-1) > 1e-12 {
		t.Fatalf("leave score = %v", leave.Score)
	}
	if leave.Score <= ScoreSelectMovie(ctx).Score || leave.Score <= ScoreContinueBrowsing(ctx).Score {
		t.Fatal("leaving should be preferred")
	}
}

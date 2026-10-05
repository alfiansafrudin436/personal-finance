package transaction

import (
	"testing"

	"personal-finance/app/transaction/repository"

	"github.com/google/uuid"
)

var (
	fromAccount = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	toAccount   = uuid.MustParse("22222222-2222-2222-2222-222222222222")
)

func TestBalanceEffects(t *testing.T) {
	cases := []struct {
		name    string
		catType repository.CategoryType
		to      uuid.NullUUID
		want    []repository.BalanceEffect
	}{
		{
			name:    "income credits the account",
			catType: repository.CategoryTypeIncome,
			want: []repository.BalanceEffect{
				{AccountID: fromAccount, Delta: "500.00"},
			},
		},
		{
			name:    "expense debits the account",
			catType: repository.CategoryTypeExpense,
			want: []repository.BalanceEffect{
				{AccountID: fromAccount, Delta: "-500.00"},
			},
		},
		{
			name:    "transfer moves the amount between accounts",
			catType: repository.CategoryTypeTransfer,
			to:      uuid.NullUUID{UUID: toAccount, Valid: true},
			want: []repository.BalanceEffect{
				{AccountID: fromAccount, Delta: "-500.00"},
				{AccountID: toAccount, Delta: "500.00"},
			},
		},
		{
			name:    "unknown type touches no balance",
			catType: repository.CategoryType("something-else"),
			want:    nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := balanceEffects(tc.catType, fromAccount, tc.to, "500.00")
			assertEffects(t, got, tc.want)
		})
	}
}

// A transfer is the case where a missing destination would corrupt a balance:
// money must never leave an account without arriving somewhere.
func TestBalanceEffectsTransferWithoutDestination(t *testing.T) {
	got := balanceEffects(repository.CategoryTypeTransfer, fromAccount, uuid.NullUUID{}, "500.00")
	if len(got) != 1 {
		t.Fatalf("got %d effects, want 1", len(got))
	}
	if got[0].Delta != "-500.00" {
		t.Errorf("delta = %q, want %q", got[0].Delta, "-500.00")
	}
}

// Reversing then reapplying is what an update does, so the two sets must
// cancel out exactly.
func TestReverseEffectsCancelsOriginal(t *testing.T) {
	original := balanceEffects(
		repository.CategoryTypeTransfer,
		fromAccount,
		uuid.NullUUID{UUID: toAccount, Valid: true},
		"1250.75",
	)

	assertEffects(t, reverseEffects(original), []repository.BalanceEffect{
		{AccountID: fromAccount, Delta: "1250.75"},
		{AccountID: toAccount, Delta: "-1250.75"},
	})
}

func assertEffects(t *testing.T, got, want []repository.BalanceEffect) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d effects, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("effect %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

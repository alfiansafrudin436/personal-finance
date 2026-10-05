package report

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"personal-finance/app/report/repository"
	"personal-finance/config"
	"personal-finance/utils"

	"github.com/labstack/echo/v4"
)

// Dashboard tuning: how much recent activity and trend history to return.
const (
	recentTransactionLimit = 5
	defaultTrendMonths     = 6
	maxTrendMonths         = 24
	topCategoryLimit       = 5
)

// Usecase handles all reporting business logic
type Usecase struct {
	repo   repository.ReportRepository
	appCfg *config.App
}

// NewUsecase creates a new report Usecase
func NewUsecase() *Usecase {
	return &Usecase{
		repo:   repository.NewRepository(),
		appCfg: config.Application,
	}
}

// GetDashboard returns everything the dashboard page needs in one round trip:
// net worth, this month income/expense/net, per-account balances, the top
// expense categories, recent activity and a monthly trend.
func (u *Usecase) GetDashboard(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	ctx := c.Request().Context()
	now := time.Now()
	monthStart := utils.StartOfMonth(now)
	monthEnd := utils.EndOfMonth(now)

	totalBalance, err := u.repo.GetTotalBalance(ctx, userID)
	if err != nil {
		log.Println("GetDashboard - GetTotalBalance error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengambil ringkasan"))
	}

	summary, err := u.repo.GetPeriodSummary(ctx, repository.GetPeriodSummaryParams{
		UserID:    userID,
		StartDate: monthStart,
		EndDate:   monthEnd,
	})
	if err != nil {
		log.Println("GetDashboard - GetPeriodSummary error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengambil ringkasan"))
	}

	accounts, err := u.repo.GetAccountBalances(ctx, userID)
	if err != nil {
		log.Println("GetDashboard - GetAccountBalances error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengambil saldo akun"))
	}

	topCategories, err := u.repo.GetSpendingByCategory(ctx, repository.GetSpendingByCategoryParams{
		UserID:    userID,
		Type:      repository.CategoryTypeExpense,
		StartDate: monthStart,
		EndDate:   monthEnd,
	})
	if err != nil {
		log.Println("GetDashboard - GetSpendingByCategory error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengambil rincian kategori"))
	}
	if len(topCategories) > topCategoryLimit {
		topCategories = topCategories[:topCategoryLimit]
	}

	recent, err := u.repo.GetRecentTransactions(ctx, repository.GetRecentTransactionsParams{
		UserID:   userID,
		RowLimit: recentTransactionLimit,
	})
	if err != nil {
		log.Println("GetDashboard - GetRecentTransactions error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengambil transaksi terbaru"))
	}

	trend, err := u.repo.GetMonthlyTrend(ctx, repository.GetMonthlyTrendParams{
		UserID:     userID,
		StartMonth: monthStart.AddDate(0, -(defaultTrendMonths - 1), 0),
		EndMonth:   monthStart,
	})
	if err != nil {
		log.Println("GetDashboard - GetMonthlyTrend error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengambil tren bulanan"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(map[string]interface{}{
		"totalBalance": totalBalance,
		"currentMonth": map[string]interface{}{
			"period":           monthStart.Format("2006-01"),
			"totalIncome":      summary.TotalIncome,
			"totalExpense":     summary.TotalExpense,
			"net":              utils.SubAmount(summary.TotalIncome, summary.TotalExpense),
			"transactionCount": summary.TransactionCount,
		},
		"accounts":           toAccountBalances(accounts),
		"topExpenseCategory": toCategoryBreakdown(topCategories),
		"recentTransactions": toRecentTransactions(recent),
		"monthlyTrend":       toTrend(trend),
	}))
}

// GetSummary returns income, expense and net for an arbitrary date range
func (u *Usecase) GetSummary(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	window, err := ParseDateRange(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
	}

	ctx := c.Request().Context()
	summary, err := u.repo.GetPeriodSummary(ctx, repository.GetPeriodSummaryParams{
		UserID:    userID,
		StartDate: window.Start,
		EndDate:   window.End,
	})
	if err != nil {
		log.Println("GetSummary - GetPeriodSummary error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengambil ringkasan"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(map[string]interface{}{
		"startDate":        window.Start.Format(utils.DateLayout),
		"endDate":          window.End.Format(utils.DateLayout),
		"totalIncome":      summary.TotalIncome,
		"totalExpense":     summary.TotalExpense,
		"net":              utils.SubAmount(summary.TotalIncome, summary.TotalExpense),
		"transactionCount": summary.TransactionCount,
	}))
}

// GetByCategory breaks a date range down per category, defaulting to expenses
func (u *Usecase) GetByCategory(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	window, err := ParseDateRange(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
	}

	reportType, valid := ParseReportType(c.QueryParam("type"))
	if !valid {
		return c.JSON(http.StatusBadRequest, utils.ResponseError("Tipe laporan tidak valid (income, expense, transfer)"))
	}

	ctx := c.Request().Context()
	rows, err := u.repo.GetSpendingByCategory(ctx, repository.GetSpendingByCategoryParams{
		UserID:    userID,
		Type:      reportType,
		StartDate: window.Start,
		EndDate:   window.End,
	})
	if err != nil {
		log.Println("GetByCategory - GetSpendingByCategory error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengambil rincian kategori"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(map[string]interface{}{
		"startDate":  window.Start.Format(utils.DateLayout),
		"endDate":    window.End.Format(utils.DateLayout),
		"type":       reportType,
		"categories": toCategoryBreakdown(rows),
	}))
}

// GetTrend returns income and expense per month for the last ?months= months,
// including months with no transactions.
func (u *Usecase) GetTrend(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	months, err := strconv.Atoi(c.QueryParam("months"))
	if err != nil || months < 1 {
		months = defaultTrendMonths
	}
	if months > maxTrendMonths {
		months = maxTrendMonths
	}

	monthStart := utils.StartOfMonth(time.Now())

	ctx := c.Request().Context()
	trend, err := u.repo.GetMonthlyTrend(ctx, repository.GetMonthlyTrendParams{
		UserID:     userID,
		StartMonth: monthStart.AddDate(0, -(months - 1), 0),
		EndMonth:   monthStart,
	})
	if err != nil {
		log.Println("GetTrend - GetMonthlyTrend error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengambil tren bulanan"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(map[string]interface{}{
		"months": months,
		"trend":  toTrend(trend),
	}))
}

// GetAccountBalances returns the balance of every active account
func (u *Usecase) GetAccountBalances(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	ctx := c.Request().Context()
	accounts, err := u.repo.GetAccountBalances(ctx, userID)
	if err != nil {
		log.Println("GetAccountBalances - GetAccountBalances error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengambil saldo akun"))
	}

	totalBalance, err := u.repo.GetTotalBalance(ctx, userID)
	if err != nil {
		log.Println("GetAccountBalances - GetTotalBalance error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengambil saldo akun"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(map[string]interface{}{
		"accounts":     toAccountBalances(accounts),
		"totalBalance": totalBalance,
	}))
}

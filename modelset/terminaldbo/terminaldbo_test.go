package terminaldbo

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"

	"github.com/tsingson/logger"

	"github.com/tsingson/goums/apis/go/goums/terminal"
	"github.com/tsingson/goums/dbv4/postgresconfig"
	"github.com/tsingson/goums/pkg/vtils"
)

var (
	cfg *postgresconfig.PostgresConfig
	log *logger.ZapLogger
)

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	log = logger.New()
	cfg = &postgresconfig.PostgresConfig{
		User:     "postgres",
		Database: "goums",
		Port:     5432,
		Host:     "127.0.0.1",
		// Host: "docker.for.mac.host.modelset",
		LogLevel: pgx.LogLevelDebug,
	}

	log = logger.New(logger.WithDebug(), logger.WithStoreInDay())

	os.Exit(m.Run())
}

func TestNewTerminalDbo(t *testing.T) {
	ctx := context.Background()
	terminalDbo, err := NewTerminalDbo(ctx, cfg, log)
	assert.NoError(t, err)
	terminalDbo.Close(ctx)
}

func TestTerminalDbo_InsertTerminal(t *testing.T) {
	ctx := context.Background()
	terminalDbo, err := NewTerminalDbo(ctx, cfg, log)
	assert.NoError(t, err)

	in := &terminal.TerminalProfileT{
		SerialNumber: vtils.RandString(16),
		ActiveCode:   vtils.RandString(16),
	}

	var id int64
	id, err = terminalDbo.Insert(ctx, in)
	assert.NoError(t, err)
	if err == nil {
		fmt.Println("id ", id)
	}
}

func TestTerminalDbo_InsertList(t *testing.T) {
	as := assert.New(t)

	ctx := context.Background()
	terminalDbo, err := NewTerminalDbo(ctx, cfg, log)
	as.NoError(err)

	list := []*terminal.TerminalProfileT{
		&terminal.TerminalProfileT{
			SerialNumber: vtils.RandString(16),
			ActiveCode:   vtils.RandString(16),
		}, &terminal.TerminalProfileT{
			SerialNumber: vtils.RandString(16),
			ActiveCode:   vtils.RandString(16),
		},
	}

	in := &terminal.TerminalListT{
		Count: int64(2),
		List:  list,
	}

	var rows int64
	rows, err = terminalDbo.InsertList(ctx, in)
	as.NoError(err)
	as.Equal(rows, in.Count)
}

func TestTerminalDbo_UpdateTerminal(t *testing.T) {
	ctx := context.Background()
	terminalDbo, err := NewTerminalDbo(ctx, cfg, log)
	assert.NoError(t, err)

	in := &terminal.TerminalProfileT{
		SerialNumber: vtils.RandString(16),
		ActiveCode:   vtils.RandString(16),
	}

	var id int64
	id, err = terminalDbo.Insert(ctx, in)
	assert.NoError(t, err)
	if err == nil {
		fmt.Println("id ", id)
	}

	in.UserID = id
	var c int64
	c, err = terminalDbo.Update(ctx, id, true, 2, 2)

	assert.NoError(t, err)
	if err == nil {
		fmt.Println(c)
	}
}

func TestTerminalDbo_Active(t *testing.T) {
	ctx := context.Background()
	terminalDbo, err := NewTerminalDbo(ctx, cfg, log)
	assert.NoError(t, err)

	in := &terminal.TerminalProfileT{
		SerialNumber: vtils.RandString(16),
		ActiveCode:   vtils.RandString(16),
	}

	var userID int64
	userID, err = terminalDbo.Insert(ctx, in)
	assert.NoError(t, err)
	apkType := "test"

	var id *terminal.TerminalProfileT
	id, err = terminalDbo.Active(ctx, in.SerialNumber, in.ActiveCode, apkType)
	assert.NoError(t, err)
	assert.Equal(t, id.UserID, userID)
}

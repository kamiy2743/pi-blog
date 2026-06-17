package appcontext

import "path/filepath"

type Context struct {
	RepoRoot string
}

func New(repoRoot string) *Context {
	return &Context{RepoRoot: repoRoot}
}

func (c *Context) Path(relativePath string) string {
	return filepath.Join(c.RepoRoot, relativePath)
}

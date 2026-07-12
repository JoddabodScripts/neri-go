package nerimity

import "testing"

func TestCacheBasicSetGet(t *testing.T) {
	c := newCache[string](0)
	c.set("a", "1")
	c.set("b", "2")
	if v, ok := c.get("a"); !ok || v != "1" {
		t.Errorf("get(a) = %q, %v, want 1, true", v, ok)
	}
	if c.len() != 2 {
		t.Errorf("len() = %d, want 2", c.len())
	}
}

func TestCacheEvictsLeastRecentlyUsed(t *testing.T) {
	c := newCache[string](2)
	c.set("a", "1")
	c.set("b", "2")
	// touch "a" so "b" becomes the least recently used.
	c.get("a")
	c.set("c", "3") // should evict "b"

	if _, ok := c.get("b"); ok {
		t.Error("expected b to be evicted")
	}
	if _, ok := c.get("a"); !ok {
		t.Error("expected a to survive eviction")
	}
	if _, ok := c.get("c"); !ok {
		t.Error("expected c to be present")
	}
	if c.len() != 2 {
		t.Errorf("len() = %d, want 2", c.len())
	}
}

func TestCacheUnboundedWhenLimitZero(t *testing.T) {
	c := newCache[int](0)
	for i := 0; i < 5000; i++ {
		c.set(string(rune(i)), i)
	}
	if c.len() != 5000 {
		t.Errorf("len() = %d, want 5000 (unbounded cache should not evict)", c.len())
	}
}

func TestCacheDelete(t *testing.T) {
	c := newCache[string](0)
	c.set("a", "1")
	c.delete("a")
	if _, ok := c.get("a"); ok {
		t.Error("expected a to be deleted")
	}
	if c.len() != 0 {
		t.Errorf("len() = %d, want 0", c.len())
	}
}

func TestCacheSetOverwritesAndPromotes(t *testing.T) {
	c := newCache[string](2)
	c.set("a", "1")
	c.set("b", "2")
	c.set("a", "1-updated") // overwrite promotes a to MRU
	c.set("c", "3")         // should evict b, not a

	if _, ok := c.get("b"); ok {
		t.Error("expected b to be evicted")
	}
	v, ok := c.get("a")
	if !ok || v != "1-updated" {
		t.Errorf("get(a) = %q, %v, want 1-updated, true", v, ok)
	}
}

package assets

import (
  "testing"
)

func TestAssetNames(t * testing.T) {
  names := AssetNames()

  for _, name := range names {
    MustAsset(name)
    info, err := AssetInfo(name)
    if err != nil {
      t.Error("Failed To Load Asset Info")
    }

    info.IsDir()
    info.Sys()
    info.ModTime()
    info.Mode()
    info.Size()
    info.Name()

    dir_names, err := AssetDir(name)
    _, _ = dir_names, err
  }
}
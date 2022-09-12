#region epilogue

import traceback

if __name__ == "__main__":
  try: Main()
  except NameError: pass
  except Exception: traceback.print_exc()
  try: main()
  except NameError: pass
  except Exception: traceback.print_exc()

#endregion epilogue

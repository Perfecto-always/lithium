package main

import (
	"syscall"
	"unsafe"

	w32 "github.com/lxn/win"
)

// Edit ID
const IDC_MAIN_EDIT = 101

func MakeIntResource(id uint16) *uint16 {
	return (*uint16)(unsafe.Pointer(uintptr(id)))
}

func WndProc(hWnd w32.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case w32.WM_CREATE:
		var hfDefault w32.HFONT
		var hEdit w32.HWND

		hEdit = w32.CreateWindowEx(w32.WS_EX_CLIENTEDGE, syscall.StringToUTF16Ptr("EDIT"), syscall.StringToUTF16Ptr(""), w32.WS_CHILD|w32.WS_VISIBLE|w32.WS_VSCROLL|w32.WS_HSCROLL|w32.ES_MULTILINE|w32.ES_AUTOVSCROLL|w32.ES_AUTOHSCROLL, 0, 0, 100, 100, hWnd, IDC_MAIN_EDIT, w32.GetModuleHandle(nil), nil)
		if hEdit == 0 {
			w32.MessageBox(hWnd, syscall.StringToUTF16Ptr("Could not create edit box."), syscall.StringToUTF16Ptr("Error"), w32.MB_OK|w32.MB_ICONERROR)
		}

		hfDefault = (w32.HFONT)(w32.GetStockObject(w32.DEFAULT_GUI_FONT))
		w32.SendMessage(hEdit, w32.WM_SETFONT, uintptr(hfDefault), lParam)
	case w32.WM_SIZE:
		var hEdit w32.HWND
		var rcClient w32.RECT

		w32.GetClientRect(hWnd, &rcClient)

		hEdit = w32.GetDlgItem(hWnd, IDC_MAIN_EDIT)
		w32.SetWindowPos(hEdit, 0, 0, 0, rcClient.Right, rcClient.Bottom, w32.SWP_NOZORDER)

	case w32.WM_CLOSE:
		w32.DestroyWindow(hWnd)
	case w32.WM_DESTROY:
		w32.PostQuitMessage(0)
	default:
		return w32.DefWindowProc(hWnd, msg, wParam, lParam)
	}
	return 0
}

func WinMain() int {
	hIcon := w32.LoadImage(0, syscall.StringToUTF16Ptr("lithium.ico"), w32.IMAGE_ICON, 0, 0, w32.LR_LOADFROMFILE)
	if hIcon == 0 {
		panic("LoadImage failed")
	}

	hInstance := w32.GetModuleHandle(nil)
	lpszClassName := syscall.StringToUTF16Ptr("WNDclass")
	var wcex w32.WNDCLASSEX

	wcex.CbSize = uint32(unsafe.Sizeof(wcex))
	wcex.Style = w32.CS_HREDRAW | w32.CS_VREDRAW
	wcex.LpfnWndProc = syscall.NewCallback(WndProc)
	wcex.CbClsExtra = 0
	wcex.CbWndExtra = 0
	wcex.HInstance = hInstance
	wcex.HIcon = w32.HICON(hIcon)
	wcex.HCursor = w32.LoadCursor(0, MakeIntResource(w32.IDC_ARROW))
	wcex.HbrBackground = w32.COLOR_WINDOW + 11

	wcex.LpszMenuName = nil

	wcex.LpszClassName = lpszClassName
	wcex.HIconSm = w32.LoadIcon(hInstance, MakeIntResource(w32.IDI_APPLICATION))
	w32.RegisterClassEx(&wcex)
	hWnd := w32.CreateWindowEx(
		w32.WS_EX_CLIENTEDGE, lpszClassName, syscall.StringToUTF16Ptr("Lithium"),
		w32.WS_OVERLAPPEDWINDOW|w32.WS_VISIBLE,
		w32.CW_USEDEFAULT, w32.CW_USEDEFAULT, 400, 400, 0, 0, hInstance, nil)

	w32.ShowWindow(hWnd, w32.SW_SHOWDEFAULT)
	w32.UpdateWindow(hWnd)

	var msg w32.MSG
	for w32.GetMessage(&msg, hWnd, 0, 0) > 0 {
		w32.TranslateMessage(&msg)
		w32.DispatchMessage(&msg)
	}

	return int(msg.WParam)
}

func main() {
	WinMain()
}

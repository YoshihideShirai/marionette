//go:build linux && marionette_desktop

package desktop

/*
#cgo pkg-config: gtk+-3.0 webkit2gtk-4.1
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>

static void marionette_destroy(GtkWidget* widget, gpointer data) {
	gtk_main_quit();
}

static void run_webkitgtk(const char* title, const char* url, int width, int height) {
	gtk_init(NULL, NULL);
	GtkWidget *window = gtk_window_new(GTK_WINDOW_TOPLEVEL);
	gtk_window_set_title(GTK_WINDOW(window), title);
	gtk_window_set_default_size(GTK_WINDOW(window), width, height);
	g_signal_connect(window, "destroy", G_CALLBACK(marionette_destroy), NULL);

	GtkWidget *webview = webkit_web_view_new();
	webkit_web_view_load_uri(WEBKIT_WEB_VIEW(webview), url);
	gtk_container_add(GTK_CONTAINER(window), webview);

	gtk_widget_show_all(window);
	gtk_main();
}
*/
import "C"
import "unsafe"

func openWebView(url string, options Options) error {
	title := C.CString(options.Title)
	defer C.free(unsafe.Pointer(title))

	targetURL := C.CString(url)
	defer C.free(unsafe.Pointer(targetURL))

	C.run_webkitgtk(title, targetURL, C.int(options.Width), C.int(options.Height))
	return nil
}

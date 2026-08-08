#import <AppKit/AppKit.h>
#import <PDFKit/PDFKit.h>
#import <WebKit/WebKit.h>
#import <dispatch/dispatch.h>

static WKWebView *qsl_find_webview(NSView *view) {
    if ([view isKindOfClass:[WKWebView class]]) {
        return (WKWebView *)view;
    }
    for (NSView *child in view.subviews) {
        WKWebView *webview = qsl_find_webview(child);
        if (webview != nil) {
            return webview;
        }
    }
    return nil;
}

static NSPrintInfo *qsl_make_print_info(double width_mm, double height_mm, int landscape, const char *printer_name) {
    NSMutableDictionary *dictionary = [[[NSPrintInfo sharedPrintInfo] dictionary] mutableCopy];
    NSPrintInfo *print_info = [[NSPrintInfo alloc] initWithDictionary:dictionary];
    CGFloat width = MAX(width_mm, 1.0) * 72.0 / 25.4;
    CGFloat height = MAX(height_mm, 1.0) * 72.0 / 25.4;
    print_info.paperSize = landscape ? NSMakeSize(height, width) : NSMakeSize(width, height);
    print_info.orientation = landscape ? NSPaperOrientationLandscape : NSPaperOrientationPortrait;
    print_info.horizontalPagination = NSClipPagination;
    print_info.verticalPagination = NSClipPagination;
    print_info.horizontallyCentered = NO;
    print_info.verticallyCentered = NO;
    print_info.topMargin = 0;
    print_info.bottomMargin = 0;
    print_info.leftMargin = 0;
    print_info.rightMargin = 0;

    if (printer_name != NULL && printer_name[0] != '\0') {
        NSString *name = [NSString stringWithUTF8String:printer_name];
        NSPrinter *printer = [NSPrinter printerWithName:name];
        if (printer != nil) {
            print_info.printer = printer;
        }
    }
    return print_info;
}

static int qsl_print_pdf_data(NSData *pdf_data, NSPrintInfo *print_info, int show_print_panel) {
    PDFDocument *document = [[PDFDocument alloc] initWithData:pdf_data];
    if (document == nil || document.pageCount == 0) {
        return 0;
    }

    NSPrintOperation *operation = [document printOperationForPrintInfo:print_info scalingMode:kPDFPrintPageScaleToFit autoRotate:NO];
    if (operation == nil) {
        return 0;
    }
    operation.showsPrintPanel = show_print_panel == 1;
    operation.showsProgressPanel = show_print_panel == 1;
    return [operation runOperation] ? 1 : 2;
}

int qsl_print_webview(double width_mm, double height_mm, int landscape, int show_print_panel, const char *printer_name) {
    if (!@available(macOS 11.0, *)) {
        return 0;
    }

    __block int status = 0;
    dispatch_semaphore_t semaphore = dispatch_semaphore_create(0);
    dispatch_block_t print_block = ^{
        NSWindow *window = [NSApp keyWindow];
        if (window == nil) {
            window = [NSApp mainWindow];
        }
        WKWebView *webview = window == nil ? nil : qsl_find_webview(window.contentView);
        if (webview == nil) {
            dispatch_semaphore_signal(semaphore);
            return;
        }

        NSPrintInfo *print_info = qsl_make_print_info(width_mm, height_mm, landscape, printer_name);
        NSString *size_script = @"(function(){var e=document.getElementById('batchPrintSheet');if(!e)return '0,0,0,0';var r=e.getBoundingClientRect();return r.left+','+r.top+','+r.width+','+e.scrollHeight;})()";
        [webview evaluateJavaScript:size_script completionHandler:^(id value, NSError *error) {
            CGFloat content_width = 1.0;
            CGFloat content_height = 1.0;
            if (error == nil && [value isKindOfClass:[NSString class]]) {
                NSArray<NSString *> *parts = [(NSString *)value componentsSeparatedByString:@","];
                if (parts.count == 4) {
                    content_width = MAX(parts[2].doubleValue, 1.0);
                    content_height = MAX(parts[3].doubleValue, 1.0);
                } else if (parts.count == 2) {
                    content_width = MAX(parts[0].doubleValue, 1.0);
                    content_height = MAX(parts[1].doubleValue, 1.0);
                }
            }

            WKPDFConfiguration *configuration = [[WKPDFConfiguration alloc] init];
            // The print sheet is laid out at the document origin in PDF mode.
            configuration.rect = CGRectMake(0, 0, content_width, content_height);
            [webview createPDFWithConfiguration:configuration completionHandler:^(NSData *pdf_data, NSError *pdf_error) {
                if (pdf_data != nil && pdf_error == nil) {
                    status = qsl_print_pdf_data(pdf_data, print_info, show_print_panel);
                }
                dispatch_semaphore_signal(semaphore);
            }];
        }];
    };

    if ([NSThread isMainThread]) {
        print_block();
        while (dispatch_semaphore_wait(semaphore, DISPATCH_TIME_NOW) != 0) {
            [[NSRunLoop currentRunLoop] runMode:NSDefaultRunLoopMode beforeDate:[NSDate dateWithTimeIntervalSinceNow:0.01]];
        }
    } else {
        dispatch_async(dispatch_get_main_queue(), print_block);
        dispatch_semaphore_wait(semaphore, DISPATCH_TIME_FOREVER);
    }
    return status;
}

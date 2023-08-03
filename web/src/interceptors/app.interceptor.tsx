import { EventTrackerService } from "@lepton-dashboard/services/event-tracker.service";
import { Injectable } from "injection-js";
import {
  HttpHandler,
  HTTPInterceptor,
  HTTPRequest,
  HTTPResponse,
} from "@lepton-dashboard/services/http-client.service";
import { catchError, mergeMap, Observable, tap, throwError } from "rxjs";
import {
  AuthService,
  UnauthorizedError,
} from "@lepton-dashboard/services/auth.service";
import { NavigateService } from "@lepton-dashboard/services/navigate.service";
import { ProfileService } from "@lepton-dashboard/services/profile.service";
import { ReactNode } from "react";
import { NotificationService } from "@lepton-dashboard/services/notification.service";
import { INTERCEPTOR_CONTEXT } from "@lepton-dashboard/interceptors/app.interceptor.context";

@Injectable()
export class AppInterceptor implements HTTPInterceptor {
  constructor(
    private navigateService: NavigateService,
    private eventTrackerService: EventTrackerService,
    private notificationService: NotificationService,
    private profileService: ProfileService,
    private authService: AuthService
  ) {}

  intercept(req: HTTPRequest, next: HttpHandler): Observable<HTTPResponse> {
    const reqHost = new URL(req.url!).host;
    const token = this.profileService.profile?.authorized_workspaces.find(
      (workspace) => new URL(workspace.auth.url).host === reqHost
    )?.auth.token;

    const headers = token
      ? { ...req.headers, Authorization: `Bearer ${token}` }
      : req.headers;

    return next
      .handle({
        ...req,
        headers,
      })
      .pipe(
        catchError((error) => {
          console.error(error);
          const status = error.status || error.response?.status;
          const ignoreErrors =
            req.context?.get(INTERCEPTOR_CONTEXT).ignoreErrors;

          const requestId = error.response?.headers?.["x-request-id"];
          const message: ReactNode = error.response?.data?.code || error.code;
          const errorMessage = error.response?.data?.message || error.message;
          const time = new Date();

          this.eventTrackerService.track("API_ERROR", {
            requestId,
            errorMessage,
            timestamp: time.toUTCString(),
          });

          const ignore401 =
            Array.isArray(ignoreErrors) && ignoreErrors.includes(401);

          // request to ignore 401 errors explicitly
          if (status === 401 && !ignore401) {
            return this.authService.logout().pipe(
              tap(() => {
                this.navigateService.navigateTo("login");
              }),
              mergeMap(() =>
                throwError(() => new UnauthorizedError("Unauthorized"))
              )
            );
          }

          if (
            ignoreErrors === true ||
            (Array.isArray(ignoreErrors) && ignoreErrors.includes(status))
          ) {
            return throwError(() => error);
          }

          const description = requestId ? (
            <>
              <strong>Error Message</strong>: {errorMessage}
              <br />
              <strong>Request ID</strong>: {requestId}
              <br />
              <strong>Timestamp</strong>: {time.toLocaleString()}
            </>
          ) : (
            errorMessage
          );

          this.notificationService.error({
            message,
            description,
          });

          return throwError(() => error);
        })
      );
  }
}

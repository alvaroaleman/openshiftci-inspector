import {JobsApi, JobsMetricsResponseBody} from "../api-client";
import INotificationService, {NotificationType} from "../notification/service/INotificationService";

export default class JobsMetricsService {

    constructor(
        private api: JobsApi,
        private notificationService: INotificationService
    ) {}

    public async getMetrics(id: string, query: string): Promise<JobsMetricsResponseBody> {
        try {
            let response = await this.api.getMetrics(id, query)
            return response.data
        } catch (e) {
            this.notificationService.notifyUser({
                type: NotificationType.ERROR,
                message: "Failed to execute query for job " + id + "."
            })
            console.log(e)
            throw e
        }
    }
}
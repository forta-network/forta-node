import 'reflect-metadata';
import { createExpressServer } from 'routing-controllers';
import { NotificationController } from './notificationController';

// creates express app, registers all controller routes and returns you express app instance
const app = createExpressServer({
  controllers: [NotificationController], // we specify controllers we want to use
});

// run express application on port 3000
app.listen(3000);
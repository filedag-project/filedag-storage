export interface RouteModel {
  name: string;
  path: string;
  component: any;
  auth: boolean;
  child?: RouteModel[];
}

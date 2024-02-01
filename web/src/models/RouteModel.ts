export interface RouteModel {
  name: string;
  path: string;
  component: any;
  auth: boolean;
  child?: RouteModel[];
}

export interface tokenType{
  isAdmin:boolean;
  parent:string;
}
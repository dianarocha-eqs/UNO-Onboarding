export enum sensorColor {
  Red = 'Red',
  Green = 'Green',
  Blue = 'Blue',
  Yellow = 'Yellow',
}

export enum sensorCategory {
  Temperature = 'Temperature',
  Humidity = 'Humidity',
  Pressure = 'Pressure',
}

export enum sensorVisibility {
  Private = 'private',
  Public = 'public',
}
  
export interface sensor {
id: number;
name: string;
color: sensorColor;
category: sensorCategory;
description: string;
visibility: sensorVisibility;
}
  
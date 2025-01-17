export interface sensor {
    id: number;
    name: string;
    color: 'Red' | 'Green' | 'Blue' | 'Yellow';
    category: 'Temperature' | 'Humidity' | 'Pressure';
    description: string,
    visibility: 'private' | 'public';
}
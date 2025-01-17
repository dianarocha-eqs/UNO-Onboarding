import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SensorAddComponent } from './sensor-add.component';

describe('SensorAddComponent', () => {
  let component: SensorAddComponent;
  let fixture: ComponentFixture<SensorAddComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SensorAddComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SensorAddComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

import {ComponentFixture, TestBed} from '@angular/core/testing';

import {DiscountDeleteComponent} from './discount-delete.component';

describe('DiscountDeleteComponent', () => {
  let component: DiscountDeleteComponent;
  let fixture: ComponentFixture<DiscountDeleteComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DiscountDeleteComponent]
    })
      .compileComponents();

    fixture = TestBed.createComponent(DiscountDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

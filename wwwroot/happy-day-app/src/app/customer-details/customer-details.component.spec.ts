import {DatePipe} from "@angular/common";
import {ComponentFixture, TestBed} from '@angular/core/testing';
import {ActivatedRoute, convertToParamMap, Router} from "@angular/router";
import {RouterTestingModule} from "@angular/router/testing";
import {ReactiveFormsModule} from "@angular/forms";
import {NoopAnimationsModule} from "@angular/platform-browser/animations";
import {MatButtonModule} from "@angular/material/button";
import {MatIconModule} from "@angular/material/icon";
import {MatInputModule} from "@angular/material/input";
import {MatFormFieldModule} from "@angular/material/form-field";

import {of} from "rxjs";
import {v4 as uuidv4} from "uuid";

import {CustomerDetailsComponent} from './customer-details.component';
import {CustomersService} from "../customers.service";
import {spyGetter} from "../tests.helper";

describe('CustomerDetailsComponent', () => {
  const customerService = jasmine.createSpyObj<CustomersService>('CustomersService', ['getById', 'update', 'create']);
  const router = jasmine.createSpyObj<Router>('Router', ['navigateByUrl'])
  const activatedRoute = jasmine.createSpyObj<ActivatedRoute>('ActivatedRoute', [], ['paramMap']);

  let component: CustomerDetailsComponent;
  let fixture: ComponentFixture<CustomerDetailsComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [CustomerDetailsComponent],
      imports: [
        NoopAnimationsModule,
        ReactiveFormsModule,
        RouterTestingModule,

        MatButtonModule,
        MatIconModule,
        MatInputModule,
        MatFormFieldModule
      ],
      providers: [
        DatePipe,
        {provide: CustomersService, useValue: customerService},
        {provide: ActivatedRoute, useValue: activatedRoute},
        {provide: Router, useValue: router},
      ]
    });
    fixture = TestBed.createComponent(CustomerDetailsComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should call getById when id is not null', () => {
    const id = uuidv4();
    spyGetter(activatedRoute, 'paramMap').and.returnValue(of(convertToParamMap({id: id})));

    customerService.getById.and.returnValue(of({
      id: id,
      name: 'name' + uuidv4(),
      comment: 'comment' + uuidv4(),
      phones: ['123456789'],
      pix: uuidv4(),
      createAt: new Date(),
      updateAt: new Date(),
    }));

    fixture.detectChanges();
    expect(customerService.getById).toHaveBeenCalledWith(id);
  });

  it('should not call getById when id is new', () => {
    spyGetter(activatedRoute, 'paramMap').and.returnValue(of(convertToParamMap({id: "new"})));

    fixture.detectChanges();
    expect(customerService.getById).not.toHaveBeenCalled();
  });

  it('should not call create/update when data is invalid', () => {
    spyGetter(activatedRoute, 'paramMap').and.returnValue(of(convertToParamMap({id: "new"})));
    fixture.detectChanges();

    component.save();
    expect(customerService.getById).not.toHaveBeenCalled();
    expect(customerService.create).not.toHaveBeenCalled();
    expect(customerService.update).not.toHaveBeenCalled();
  });

  it('should have required error on name when name is empty', () => {
    spyGetter(activatedRoute, 'paramMap').and.returnValue(of(convertToParamMap({id: "new"})));
    fixture.detectChanges();

    component.form.patchValue({name: ''});
    expect(component.form.get('name')!.hasError('required')).toBeTruthy();
  });

  it('should have max length error on name when name length is large than 255', () => {
    spyGetter(activatedRoute, 'paramMap').and.returnValue(of(convertToParamMap({id: "new"})));
    fixture.detectChanges();

    component.form.patchValue({name: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Sollicitudin aliquam ultrices sagittis orci a scelerisque purus semper eget. Arcu ac tortor dignissim convallis aenean. At quis risus sed vulputate odio ut.'});
    expect(component.form.get('name')!.hasError('maxlength')).toBeTruthy();
  });

  it('should have max length error on pix when pix length is large than 255', () => {
    spyGetter(activatedRoute, 'paramMap').and.returnValue(of(convertToParamMap({id: "new"})));
    fixture.detectChanges();

    component.form.patchValue({
      name: 'test',
      pix: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Sollicitudin aliquam ultrices sagittis orci a scelerisque purus semper eget. Arcu ac tortor dignissim convallis aenean. At quis risus sed vulputate odio ut.'
    });
    expect(component.form.get('pix')!.hasError('maxlength')).toBeTruthy();
  });

  it('should have required error on phone when phone is empty', () => {
    spyGetter(activatedRoute, 'paramMap').and.returnValue(of(convertToParamMap({id: "new"})));
    fixture.detectChanges();

    component.addPhone('')
    expect(component.phones.at(0).hasError('required')).toBeTruthy();
  });

  it('should have min length error on phone when phone length is less than 8', () => {
    spyGetter(activatedRoute, 'paramMap').and.returnValue(of(convertToParamMap({id: "new"})));
    fixture.detectChanges();

    component.addPhone('1234567');
    expect(component.phones.at(0).hasError('minlength')).toBeTruthy();
  });

  it('should have max length error on phone when phone length is large than 11', () => {
    spyGetter(activatedRoute, 'paramMap').and.returnValue(of(convertToParamMap({id: "new"})));
    fixture.detectChanges();

    component.addPhone('123456789012');
    expect(component.phones.at(0).hasError('maxlength')).toBeTruthy();
  });

  it('should have pattern error on phone when phone length is large than 11', () => {
    spyGetter(activatedRoute, 'paramMap').and.returnValue(of(convertToParamMap({id: "new"})));
    fixture.detectChanges();

    component.addPhone('abcdefgh');
    expect(component.phones.at(0).hasError('pattern')).toBeTruthy();
  });

  it('should call create when id is new', () => {
    spyGetter(activatedRoute, 'paramMap').and.returnValue(of(convertToParamMap({id: "new"})));
    customerService.create.and.returnValue(of({
      id: uuidv4(),
      name: 'name',
      comment: 'comment',
      phones: ['123456789'],
      pix: 'pix',
      createAt: new Date(),
      updateAt: new Date(),
    }));

    fixture.detectChanges();

    component.addPhone('123456789');
    component.form.patchValue({
      name: 'name',
      comment: 'comment',
      pix: 'pix',
    });

    component.save();
    expect(customerService.create).toHaveBeenCalled();
  });

  it('should call update when id is not null', () => {
    const id = uuidv4();
    spyGetter(activatedRoute, 'paramMap').and.returnValue(of(convertToParamMap({id: id})));
    customerService.update.and.returnValue(of({
      id: id,
      name: 'name',
      comment: 'comment',
      phones: ['123456789'],
      pix: 'pix',
      createAt: new Date(),
      updateAt: new Date(),
    }));

    fixture.detectChanges();

    component.addPhone('123456789');
    component.form.patchValue({
      name: 'name',
      comment: 'comment',
      pix: 'pix',
    });

    component.save();
    expect(customerService.update).toHaveBeenCalled();
  });

  it('should call navigateByUrl when cancel is called', () => {
    spyGetter(activatedRoute, 'paramMap').and.returnValue(of(convertToParamMap({id: "new"})));
    router.navigateByUrl.and.returnValue(Promise.resolve(true));

    fixture.detectChanges();

    component.cancel();
    expect(router.navigateByUrl).toHaveBeenCalledWith('/customers');
  });
});

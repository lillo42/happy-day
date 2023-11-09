import {inject, TestBed} from '@angular/core/testing';

import { CustomersService } from './customers.service';
import { HttpClientTestingModule, HttpTestingController } from "@angular/common/http/testing";

import { v4 as uuidv4 } from "uuid";

describe('CustomersService', () => {

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [ HttpClientTestingModule ],
      providers: [ CustomersService ]
    });
  });

  it('should call delete with the given id', () => {
    inject([HttpTestingController, CustomersService],
      (httpMock: HttpTestingController ,service: CustomersService) => {
        const id = uuidv4();
        service.delete(id).subscribe();
        const req = httpMock.expectOne(`/api/customers/${id}`);
        expect(req.request.method).toEqual('DELETE');
        req.flush({});
        httpMock.verify();
    });
  });

  it('should call update with the given id and customer', () => {
    inject([HttpTestingController, CustomersService],
      (httpMock: HttpTestingController ,service: CustomersService) => {
        const id = uuidv4();
        const customer = {
          id: id,
          name:`test-name-${id}}`,
          comment: `test-comment-${id}}`,
          phones: [`test-phone-${id}}`],
          pix: '' ,
          createAt: new Date(),
          updateAt: new Date()
        };
        service.update(id, customer).subscribe();
        const req = httpMock.expectOne(`/api/customers/${id}`);
        expect(req.request.method).toEqual('PUT');
        req.flush({});
        httpMock.verify();
    });
  });

  it('should call create with the given customer', () => {
    inject([HttpTestingController, CustomersService],
      (httpMock: HttpTestingController ,service: CustomersService) => {
        const id = uuidv4();
        const customer = {
          id: '',
          name:`test-name-${id}}`,
          comment: `test-comment-${id}}`,
          phones: [`test-phone-${id}}`],
          pix: '',
          createAt: new Date(),
          updateAt: new Date()
        };
        service.create(customer).subscribe();
        const req = httpMock.expectOne(`/api/customers/`);
        expect(req.request.method).toEqual('POST');
        req.flush({});
        httpMock.verify();
    });
  });

  it('should call getById with the given id', () => {
    inject([HttpTestingController, CustomersService],
      (httpMock: HttpTestingController ,service: CustomersService) => {
        const id = uuidv4();
        service.getById(id).subscribe();
        const req = httpMock.expectOne(`/api/customers/${id}`);
        expect(req.request.method).toEqual('GET');
        req.flush({});
        httpMock.verify();
    });
  });

  it('should call get with the given params', () => {
    inject([HttpTestingController, CustomersService],
      (httpMock: HttpTestingController ,service: CustomersService) => {
        const name = 'test-name';
        const phone = 'test-phone';
        const comment = 'test-comment';
        const page = 1;
        const size = 10;
        service.get(name, phone, comment, page, size).subscribe();
        const req = httpMock.expectOne(`/api/customers?page=${page}&size=${size}&name=${name}&phone=${phone}&comment=${comment}`);
        expect(req.request.method).toEqual('GET');
        req.flush({});
        httpMock.verify();
    });
  });
});

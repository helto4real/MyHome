import { Action, ActionCreator } from 'redux';
import { ThunkAction } from 'redux-thunk';
import { EntityState } from '../reducers/entity.js';
import { ShopActionGetProducts } from './shop.js';

export const STATE_UPDATE = 'STATE_UPDATE';

export interface EnityActionUpdateState extends Action<'STATE_UPDATE'> {products: EntityState};

export type EntityAction = EnityActionUpdateState; //| ShopActionAddToCart | ShopActionRemoveFromCart | ShopActionCheckoutSuccess | ShopActionCheckoutFailure;

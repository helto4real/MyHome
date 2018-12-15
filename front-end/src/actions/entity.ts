import { Action, ActionCreator } from 'redux';
import { ThunkAction } from 'redux-thunk';
import { EntitiesState } from '../reducers/entity.js';
//import { ShopActionGetProducts } from './shop.js';
import { RootState } from '../store.js';

export const STATE_UPDATE = 'STATE_UPDATE';
export const GET_ENTITIES = 'GET_ENTITIES';

export interface EntityActionGetEntities extends Action<'GET_ENTITIES'> {entities: EntitiesState};
export interface EntityActionUpdateState extends Action<'STATE_UPDATE'> {entityId: string, state: string};

export type EntityAction = EntityActionGetEntities | EntityActionUpdateState; 

type ThunkResult = ThunkAction<void, RootState, undefined, EntityAction>;

const ENTITY_LIST = [
    {"id": "switch.switch1", "state": "on", "attributes": "{'test': 'test'}"},
    {"id": "switch.switch2", "state": "on", "attributes": "{'test': 'test'}"},
    {"id": "light.light1", "state": "off", "attributes": "{'test': 'test'}"},
    {"id": "light.light2", "state": "on", "attributes": "{'test': 'test'}"},
];

export const getAllEntities: ActionCreator<ThunkResult> = () => (dispatch) => {
    // Here you would normally get the data from the server. We're simulating
    // that by dispatching an async action (that you would dispatch when you
    // succesfully got the data back)
  
    // You could reformat the data in the right format as well:
    const entities = ENTITY_LIST.reduce((obj, entity) => {
      obj[entity.id] = entity
      return obj
    }, {} as EntitiesState);
  
    dispatch({
      type: GET_ENTITIES,
      entities
    });
  };
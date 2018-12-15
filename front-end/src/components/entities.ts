import { LitElement, html, property } from '@polymer/lit-element';
import { connect } from 'pwa-helpers/connect-mixin.js';

// This element is connected to the Redux store.
import { store, RootState } from '../store.js';
import { ButtonSharedStyles } from './button-shared-styles.js';
import { EntitiesState } from '../reducers/entity.js';
import { getAllEntities } from '../actions/entity.js';

import my_home from '../reducers/entity.js';
store.addReducers({
  my_home
});

class EntityList extends connect(store)(LitElement) {
    protected render() {
      return html`
        ${ButtonSharedStyles}
        <style>
          :host { display: block; }
        </style>
        <p ?hidden="${Object.keys(this._entities).length !== 0}">No entities...</p>
        <ul>
        ${Object.keys(this._entities).map((key) => {
        const item = this._entities[key];
        return html`
          <li>${item.id}</li>
          `
      })}
        </ul>
      `;
    }
  
    protected firstUpdated() {
        store.dispatch(getAllEntities());
    }

 
    @property({type: Object})
    private _entities: EntitiesState = {};

    // This is called every time something is updated in the store.
    //stateChanged(state: RootState) {
    // This is called every time something is updated in the store.
    stateChanged(state: RootState) {
      if (state.my_home) {
        this._entities = state.my_home!.entities;
      }
      else {
        console.error("state entity is null!")
        console.error(state);
      }
      
    }
    
}

window.customElements.define('entity-list', EntityList);

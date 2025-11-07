

## `window.pathless`

The `window.pathless` object provides the API facilitating interaction between `panels` (viewport) and `frames` (html/js/css).    

#### `pathless.context()`
Returns the DOM element of the focused panel, its current frame, and its panel specific state.

#### `pathless.update(key, value)`

Automatically detects the focused panel and current frame, uses key-value pair's to save and restore state across navigation and layout changes.

#### `pathless.fetch(url, opts)`
Fetches data from `url` (string) with standard fetch `opts` (object, optional). Returns a promise resolving to `{ data, headers }`. 

- `data`: Parsed response (blob URL for images, JSON for JSON responses, text otherwise)
- `headers`: Response headers object

Caching and request deduplication available via `opts.key`. 

#### `pathless.onKey(handler)`

Event handler used to register keybinds in frames, automatically scoped to the focused panel.


### pathless example

```javascript
const { panel, state } = pathless.context();
const key = 'slideIndex';
    
   let slides = [];
   let index = state[key] || 0;

   async function show(i) {
      if (!slides.length) return;
        
      index = ((i %% slides.length) + slides.length) %% slides.length;
      pathless.update(key, index);

      const img = panel.querySelector('.slides img');
      if (!img) return;
        
      const slide = slides[index];
      const { data } = await pathless.fetch(apiUrl + '/%s/' + slide, { key: slide });
      img.src = data;
      img.alt = slide;
    }

    pathless.fetch(apiUrl + '/%s/order', { key: 'order-%s' })
          .then(({ data }) => {
              slides = data;
              if (slides.length) show(index);
          });

    pathless.onKey((k) => {
        if (k === 'a') show(index - 1);
        else if (k === 'd') show(index + 1);
    });
```

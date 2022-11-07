using DevSmtp.Core.Stores;

namespace DevSmtp.Core.Commands
{
    public sealed class TurnHandler : ICommandHandler<Turn, TurnResult>
    {
        private readonly IDataStore _dataStore;

        public TurnHandler(IDataStore dataStore)
        {
            this._dataStore = dataStore ?? throw new ArgumentNullException(nameof(dataStore));
        }

        public Task<TurnResult> ExecuteAsync(Turn command, CancellationToken cancellationToken = default)
        {
            throw new NotImplementedException();
        }
    }
}

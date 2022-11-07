using DevSmtp.Core.Stores;

namespace DevSmtp.Core.Commands
{
    public sealed class RsetHandler : ICommandHandler<Rset, RsetResult>
    {
        private readonly IDataStore _dataStore;

        public RsetHandler(IDataStore dataStore)
        {
            this._dataStore = dataStore ?? throw new ArgumentNullException(nameof(dataStore));
        }

        public Task<RsetResult> ExecuteAsync(Rset command, CancellationToken cancellationToken = default)
        {
            throw new NotImplementedException();
        }
    }
}
